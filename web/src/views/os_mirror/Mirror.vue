<i18n>
en:
  titleRepo: OS mirrors
  repoDescription1: You muse provide a mirror address (e.g. CentOS yum mirror, Ubuntu apt mirror) that all the nodes in your mirror could access.
  repoDescription2: Usually, a enterprise has its own os mirror, if not, Kuboard alos provide a wizard, so that you can create an os mirror quickly.
  url: Mirror/Repo Url
  url_placeholder: Nodes use this url to access the OS Mirror
  created: Created
  success: Ready
  failed: Not Ready
  basics: Base information
zh:
  titleRepo: OS 软件源
  repoDescription1: 您必须定义您集群机器可以访问的操作系统软件源地址（例如 CentOS 的 yum 源、Ubuntu 的 apt 源等）以供使用
  repoDescription2: 通常企业都有自己的常用操作系统的本地软件源，Kuboard 提供向导，帮助您快速设置一个操作系统软件源
  provision: 部署软件源
  url: 软件源地址
  url_placeholder: 各 K8S 节点访问此软件源时所使用的 URL 地址
  created: 已创建
  success: 可用
  failed: 不可用
  basics: 基本信息
</i18n>

<template>
  <div>
    <ControlBar :title="name">
      <template v-if="os_mirror">
        <template v-if="mode === 'view'">
          <el-button v-if="!os_mirror.history.processing" type="primary" icon="el-icon-edit" @click="$router.push(`${name}?mode=edit`)">{{ $t('msg.edit') }}</el-button>
          <MirrorProcessing v-if="os_mirror.status.kind === 'provision'" :os_mirror="os_mirror" :loading="loading" :name="name" @refresh="refresh"></MirrorProcessing>
        </template>
        <template v-else>
          <el-popconfirm :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')" placement="bottom-start"
            icon="el-icon-warning" icon-color="red" :title="$t('msg.confirmToCancel')" @confirm="cancelEdit">
            <template #reference>
              <el-button type="default" icon="el-icon-close">{{$t('msg.cancel')}}</el-button>
            </template>
          </el-popconfirm>
          <el-button type="primary" icon="el-icon-check" :disabled="noSaveRequired" @click="save">{{$t('msg.save')}}</el-button>
        </template>
      </template>
    </ControlBar>
    <el-form v-if="os_mirror" ref="form" :model="os_mirror" label-position="left" label-width="150px" @submit.prevent.stop>
      <el-card shadow="never" v-if="loading">
        <el-skeleton animated :rows="10"></el-skeleton>
      </el-card>
      <div v-else-if="os_mirror.status">
        <el-card class="app_margin_bottom" shadow="never">
          <div style="margin-bottom: -20px;">
            <el-form-item :label="t('url')" prop="status.url" :rules="urlRules">
              <span v-if="mode === 'view'">
                <KuboardSprayLink :href="os_mirror.status.url">{{os_mirror.status.url}}</KuboardSprayLink>
              </span>
              <el-input v-else v-model.trim="mirror_url" :placeholder="t('url_placeholder')"></el-input>
            </el-form-item>
            <el-form-item :label="$t('msg.status')">
              <el-tag v-if="os_mirror.status.status === 'created'" type="">{{t('created')}}</el-tag>
              <el-tag v-if="os_mirror.status.status === 'success'" type="success" effect="dark">{{t('success')}}</el-tag>
              <el-tag v-if="os_mirror.status.status === 'failed'" type="danger" effect="dark">{{t('failed')}}</el-tag>
            </el-form-item>
          </div>
        </el-card>
        <el-tabs type="border-card" v-model="currentTab">
          <el-tab-pane :label="t('basics')" name="basics">
            <Params v-if="os_mirror" :os_mirror="os_mirror"></Params>
          </el-tab-pane>
          <el-tab-pane v-if="os_mirror.status.kind === 'provision'" :label="t('provision')" name="provision">
            <el-scrollbar height="calc(100vh - 345px)">
              <Provision v-if="os_mirror.inventory" :os_mirror="os_mirror"></Provision>
            </el-scrollbar>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-form>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import {computed} from 'vue'
import Provision from './provision/Provision.vue'
import Params from './params/Params.vue'
import MirrorProcessing from './MirrorProcessing.vue'

export default {
  mixins: [mixin],
  props: {
    name: { type: String, required: true },
    mode: { type: String, required: false, default: 'view' }
  },
  data() {
    return {
      percentage: 0,
      os_mirror: undefined,
      original_data: '',
      loading: false,
      currentTab: 'basics',
      urlRules: [
        { required: true, type: 'string', message: 'url is required.', trigger: 'blur' },
      ]
    }
  },
  provide () {
    return {
      editMode: computed(() => {
        if (this.mode === 'view') {
          return 'view'
        }
        if (this.os_mirror && this.os_mirror.history.success_tasks.length > 0) {
          return 'frozen'
        }
        return this.mode
      }),
      isInstalled: computed(() => {
        if (this.os_mirror === undefined) {
          return false
        }
        if (this.os_mirror.history.success_tasks.length > 0) {
          return true
        }
        return false
      })
    }
  },
  refresh () {
    this.refresh()
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: this.t('titleRepo'), to: '/settings/mirrors' },
      { label: this.name },
    ]
  },
  computed: {
    mirror_url: {
      get () {
        return this.os_mirror.status.url
      },
      set (v) {
        this.os_mirror.status.url = v
      }
    },
    noSaveRequired () {
      return this.original_data === JSON.stringify(this.os_mirror)
    }
  },
  components: { Provision, MirrorProcessing, Params },
  mounted () {
    this.refresh()
  },
  methods: {
    async refresh () {
      this.loading = true
      await this.kuboardSprayApi.get(`/mirrors/${this.name}`).then(resp => {
        this.os_mirror = resp.data.data
        this.original_data = JSON.stringify(this.os_mirror)
      }).catch(e => {
        console.log(e)
      })
      this.loading = false
    },
    cancelEdit () {
      this.$router.replace(this.name)
      this.refresh()
    },
    save () {
      this.$refs.form.validate(flag => {
        if (flag) {
          this.kuboardSprayApi.put(`/mirrors/${this.name}`, this.os_mirror).then(() => {
            this.$message.success(this.$t('msg.save_succeeded'))
            this.refresh()
            this.$router.replace(`${this.name}`)
          }).catch(e => {
            this.$message.error(this.$t('msg.save_failed', e.response.data.message))
          })
        }
      })
    }
  }
}
</script>

<style scoped lang="css">

</style>
