<i18n>
en:
  bastionUsage: KuboardSpray can access Kubernetes Cluster Nodes through bastion.
  addSshKey: Add Private Key
  ansible_host_placeholder: 'KuboardSpray use this ip or hostname to connect to the node.'
  default_value: 'Default: {default_value} (inhirit from value configured in Global Config tab)'
  duplicateIP: "IP address conflict with {node}"
  description: Access target nodes throguh bastion/jumpserver
zh:
  bastionUsage: KuboardSpray 可以通过跳板机或堡垒机访问将要安装 K8S 集群的目标节点。
  addSshKey: 添加私钥
  ansible_host_placeholder: 'KuboardSpray 连接该主机时所使用的主机名或 IP 地址'
  default_value: '默认值：{default_value} （继承自全局设置标签页中的配置）'
  duplicateIP: "IP 地址不能与其他节点相同：{node}"
  description: 通过跳板机或者堡垒机访问目标节点
  protocol: 跳板机协议
</i18n>


<template>
  <ConfigSection ref="configSection" v-model:enabled="bastionEnabled" :label="$t('obj.bastion')"
    :description="t('description')" anti-freeze>
    <template #more>
      {{ t('bastionUsage') }}
    </template>
    <el-form-item :label="t('protocol')" required>
      <el-radio-group v-model="computedBastionType" :disabled="editMode != 'edit'">
        <el-radio value="ssh">SSH</el-radio>
        <el-radio value="socks5">SOCKS5</el-radio>
      </el-radio-group>
    </el-form-item>
    <FieldString :holder="holder" fieldName="ansible_host" :prop="`all.hosts.${nodeName}`" anti-freeze
      :placeholder="t('ansible_host_placeholder')" :rules="hostRules"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_port" :prop="`all.hosts.${nodeName}`"
      :placeholder="placeholder('ansible_port')" anti-freeze required></FieldString>
    <template v-if="computedBastionType == 'ssh'">
      <FieldString :holder="holder" fieldName="ansible_user" :prop="`all.hosts.${nodeName}`"
        :placeholder="placeholder('ansible_user')" anti-freeze required></FieldString>
      <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList" anti-freeze
        clearable :placeholder="placeholder('ansible_ssh_private_key_file')">
        <template #edit>
          <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus"
            @click="$refs.addPrivateKey.show()">{{ t('addSshKey') }}</el-button>
        </template>
      </FieldSelect>
      <FieldString :holder="holder" fieldName="ansible_password" show-password anti-freeze clearable
        :placeholder="placeholder('ansible_password')"></FieldString>
      <SshAddPrivateKey ref="addPrivateKey" ownerType="cluster" :ownerName="cluster.name"></SshAddPrivateKey>
    </template>
    <template v-else-if="computedBastionType == 'socks5'">
      <FieldString :holder="holder" fieldName="ansible_user" :prop="`all.hosts.${nodeName}`"
        :placeholder="placeholder('ansible_user')" anti-freeze></FieldString>
      <FieldString :holder="holder" fieldName="ansible_password" show-password anti-freeze clearable
        :placeholder="placeholder('ansible_password')"></FieldString>
    </template>
  </ConfigSection>
</template>

<script>
import SshAddPrivateKey from '../../../private_key/SshAddPrivateKey.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: false },
    holder: { type: Object, required: true },
    prop: { type: String, required: true },
  },
  data() {
    return {
      hostRules: [
        { required: true, message: this.$t('field.ansible_host') + this.$t('field.is_required_field'), trigger: 'blur' },
      ]
    }
  },
  inject: ['editMode'],
  computed: {
    computedBastionType: {
      get() {
        if (this.inventory.all.hosts.bastion && this.inventory.all.hosts.bastion.bastion_type) {
          return this.inventory.all.hosts.bastion.bastion_type
        }
        return "ssh"
      },
      set(type) {
        this.inventory.all.hosts.bastion.bastion_type = type;
      }
    },
    inventory: {
      get() { return this.cluster.inventory },
      set() { }
    },
    bastionEnabled: {
      get() {
        return this.inventory.all.hosts.bastion !== undefined
      },
      set(v) {
        if (v) {
          this.inventory.all.hosts.bastion = this.inventory.all.hosts.bastion || { ansible_host: '', ansible_user: '' }
        } else {
          delete this.inventory.all.hosts.bastion
          delete this.inventory.all.children.target.vars.ansible_ssh_common_args
        }
      }
    },
    holderRef: {
      get() { return this.holder },
      set() { }
    },
    ansible_become: {
      get() {
        if (this.holder.ansible_become !== undefined) {
          return this.holder.ansible_become
        }
        return this.cluster.inventory.all.children.target.vars.ansible_become
      },
      set(v) {
        this.holderRef.ansible_become = v
      }
    }
  },
  components: { SshAddPrivateKey },
  mounted() {
  },
  watch: {
    holder: {
      handler: function (bastion) {
        if (bastion && bastion.ansible_host) {
          if (this.computedBastionType == 'ssh') {
            let sshPass = ''
            if (bastion['ansible_password']) {
              sshPass = `sshpass -p '${bastion['ansible_password']}' `
            }
            let temp = `-o ProxyCommand="${sshPass}ssh -F /dev/null -o ControlMaster=auto -o ControlPersist=30m -o ControlPath={{kuboardspray_cluster_dir}}/%%r@%%h:%%p -o ConnectTimeout=10 -o ConnectionAttempts=100 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -W %h:%p -p`
            temp += bastion["ansible_port"] + " " + bastion["ansible_user"] + "@" + bastion["ansible_host"]
            if (bastion["ansible_ssh_private_key_file"]) {
              temp += " -i " + bastion["ansible_ssh_private_key_file"]
            }
            temp += '"'
            this.inventory.all.children.target.vars.ansible_ssh_common_args = temp
          } else if (this.computedBastionType == 'socks5') {
            let temp = `-o ProxyCommand="nc -X 5 -x ${bastion.ansible_host}:${bastion.ansible_port} %h %p"`
            if (bastion['ansible_user']) {
              temp = `-o ProxyCommand="nc -X 5 -x ${bastion['ansible_user']}@${bastion.ansible_host}:${bastion.ansible_port} %h %p"`
            }
            if (bastion['ansible_password'] && bastion['ansible_user']) {
              temp = `-o ProxyCommand="nc -X 5 -x ${bastion['ansible_user']}:${bastion['ansible_password']}@${bastion.ansible_host}:${bastion.ansible_port} %h %p"`
            }
            this.inventory.all.children.target.vars.ansible_ssh_common_args = temp
          } else {
            delete this.inventory.all.children.target.vars.ansible_ssh_common_args
          }
        } else {
          delete this.inventory.all.children.target.vars.ansible_ssh_common_args
        }
      },
      deep: true,
    }
  },
  methods: {
    placeholder(fieldName) {
      return this.$t('field.' + fieldName + '_placeholder')
    },
    async loadSshKeyList() {
      let result = []
      await this.kuboardSprayApi.get(`/private-keys/cluster/${this.cluster.name}`).then(resp => {
        for (let item of resp.data.data) {
          result.push({ label: item, value: '{{ kuboardspray_cluster_dir }}/private-key/' + item })
        }
      })
      return result
    },
  }
}
</script>

<style scoped lang="css"></style>
