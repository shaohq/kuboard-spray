<i18n>
en:
  label: Metrics
  description: metrics_server
  addon_function: Regularly takes metris info from K8S components, a must for features like Horizontal Auto Scaler.
  metrics_server_version: Version
zh:
  label: Metrics
  description: metrics_server
  addon_function: 定时采集 K8S 各组件的性能数据，是容器组水平伸缩等功能的基础依赖。
  metrics_server_version: 版本
</i18n>

<template>
  <AddonSection v-model:enabled="enabled" :label="t('label')" :description="$t('obj.addon', {name: this.t('description')})"
    :cluster="cluster" addonName="metrics_server" @refresh="$emit('refresh')">
    <template #more>
      {{t('addon_function')}}
    </template>
    <FieldString :holder="vars" :prop="prop" fieldName="metrics_server_metric_resolution" :rules="resolutionRules"></FieldString>
    <FieldString :holder="addon.params" fieldName="metrics_server_version" :label="t('metrics_server_version')" readOnly></FieldString>
  </AddonSection>
</template>

<script>
import AddonSection from '../AddonSection.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    addon: { type: Object, required: true },
  },
  data() {
    return {
      resolutionRules: [
        {
          validator: (rule, value, callback) => {
            if (value) {
              if (/^[0-9]+s/.test(value)) {
                return callback()
              } else {
                return callback('必须是数字，并加上 s 作为结尾')
              }
            }
            return callback()
          },
          trigger: 'blur'
        }
      ]
    }
  },
  computed: {
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {},
    },
    prop () {
      return 'all.children.target.children.k8s_cluster.vars'
    },
    enabled: {
      get () {
        return this.vars.metrics_server_enabled
      },
      set (v) {
        this.vars.metrics_server_enabled = v
      }
    }
  },
  components: { AddonSection },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
