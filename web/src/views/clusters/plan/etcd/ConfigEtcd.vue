<i18n>
en:
  ETCD: ETCD Params (Apply to all the etcd nodes)
zh:
  ETCD: ETCD 参数（适用范围：所有 etcd 节点）
</i18n>

<template>
  <div>
    <ConfigSection v-model:enabled="enabledEtcd" :description="t('ETCD')" label="ETCD" disabled anti-freeze>
      <FieldSelect :holder="cluster.inventory.all.children.target.children.etcd.vars" prop="all.children.target.children.etcd.vars" required
        fieldName="etcd_deployment_type" :loadOptions="loadEtcdDeploymentOptions"></FieldSelect>
      <!-- <el-form-item :label="$t('field.etcd_deployment_type')">
        <span class="app_text_mono">{{etcd_deployment_type}}</span>
      </el-form-item> -->
      <FieldString :holder="cluster.inventory.all.children.target.children.etcd.vars" prop="all.children.target.children.etcd.vars" disabled
        fieldName="etcd_data_dir"></FieldString>
    </ConfigSection>
  </div>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    enabledEtcd: {
      get () {
        return true
      },
      set (v) {
        console.log(v)
      }
    },
    etcd_deployment_type () {
      if (this.cluster) {
        return this.$t('field.etcd_deployment_type-' + this.cluster.inventory.all.children.target.children.etcd.vars.etcd_deployment_type)
      } else {
        return ''
      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    async loadEtcdDeploymentOptions () {
      let result = []
      for (let item of this.cluster.resourcePackage.data.etcd.etcd_deployment_type) {
        result.push({
          label: this.$t('field.etcd_deployment_type-' + item),
          value: item,
          disabled: item === 'docker' ? this.cluster.inventory.all.children.target.vars.container_manager !== 'docker' : false,
        })
      }
      return result
    },
  }
}
</script>

<style scoped lang="css">

</style>
