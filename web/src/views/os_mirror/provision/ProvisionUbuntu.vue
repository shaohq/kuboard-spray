<i18n>
en:
  special_time: Frequncy to sync
  KB: KiloBytes/Thread
  kb: about { kb } Mb/s per thread
  threadCount: Thread Count
  totalkb: about { kb } Mb/s in total
  releases: Releases included
zh:
  special_time: 从上游同步的频率
  KB: 千字节每秒每线程
  kb: 大约 { kb } Mb/s 每线程
  threadCount: 个下载线程
  totalkb: 总共大约 { kb } Mb/s
  release: 包含的版本
</i18n>

<template>
  <div>
    <!-- <FieldNumber :holder="vars" :prop="prop" fieldName="apache2_default_port" required></FieldNumber> -->
    <FieldString :holder="vars" :prop="prop" fieldName="apt_mirror_dir" required disabled></FieldString>
    <FieldCommon required :holder="vars" fieldName="apt_mirror_ubuntu_mirror" :prop="prop">
      <template #edit>
        <el-select v-model="apt_mirror_ubuntu_mirror" style="width: 100%;" :placeholder="$t('field.apt_mirror_ubuntu_mirror')" allow-create filterable>
          <el-option v-for="(url, index) in mirrorOptions[os_mirror.status.type]" :key="'url' + index" :value="url">
            {{ url }}
          </el-option>
        </el-select>
      </template>
      <template #view>
        <span class="app_text_mono">{{apt_mirror_ubuntu_mirror}}</span>
      </template>
    </FieldCommon>
    <FieldCommon :label="t('release')" required anti-freeze :holder="vars" fieldName="apt_mirror_repos" :prop="prop">
      <template #edit>
        <el-checkbox-group v-model="apt_mirror_repos">
          <el-checkbox v-for="(value, key) in releases" :key="key + 'release'" :label="key" :value="key">
            <span class="app_text_mono">
              Ubuntu {{value}} - {{key}}
            </span>
          </el-checkbox>
        </el-checkbox-group>
      </template>
      <template #view>
        <el-checkbox-group v-model="apt_mirror_repos" disabled>
          <el-checkbox v-for="(value, key) in releases" :key="key + 'release'" :label="key" :value="key">
            <span class="app_text_mono">
              Ubuntu {{value}} - {{key}}
            </span>
          </el-checkbox>
        </el-checkbox-group>
      </template>
    </FieldCommon>
    <FieldCommon :label="t('architecture')" required anti-freeze>
      <template #edit>
        <el-checkbox-group  v-model="apt_mirror_architecture">
          <el-checkbox v-for="(value, key) in architecture[os_mirror.status.type]" :key="key + 'release'" :label="key" :value="key">
            <span class="app_text_mono">
              {{key}}
            </span>
          </el-checkbox>
        </el-checkbox-group>
      </template>
      <template #view>
        <el-checkbox-group v-model="apt_mirror_architecture" disabled>
          <el-checkbox v-for="(value, key) in architecture[os_mirror.status.type]" :key="key + 'release'" :label="key" :value="key">
            <span class="app_text_mono">
              {{key}}
            </span>
          </el-checkbox>
        </el-checkbox-group>
      </template>
    </FieldCommon>
    <FieldBool :holder="vars" :prop="prop" fieldName="apt_mirror_schedule_updates" anti-freeze></FieldBool>
    <FieldRadio v-if="vars.apt_mirror_schedule"
      :disabled="!vars.apt_mirror_schedule_updates"
      :holder="vars.apt_mirror_schedule[0]"
      fieldName="special_time" :label="t('special_time')"
      :options="['daily', 'weekly', 'monthly']"></FieldRadio>
    <FieldBool :holder="vars" :prop="prop" fieldName="apt_mirror_populate_repos" anti-freeze></FieldBool>
    <FieldBool :holder="vars" :prop="prop" fieldName="apt_mirror_enable_limit_rate" anti-freeze></FieldBool>
    <FieldNumber :holder="vars" :prop="prop" fieldName="apt_mirror_limit_rate" anti-freeze
      v-if="vars.apt_mirror_enable_limit_rate">
      <template #append>{{ t('KB') }}</template>
      <div class="placeholder">{{ t('kb', { kb: vars.apt_mirror_limit_rate * 8 / 1000 }) }}</div>
    </FieldNumber>
    <FieldNumber :holder="vars" :prop="prop" fieldName="apt_mirror_nthreads" anti-freeze
      v-if="vars.apt_mirror_enable_limit_rate">
      <template #append>{{ t('threadCount') }}</template>
      <div class="placeholder">{{ t('totalkb', { kb: vars.apt_mirror_limit_rate * 8 * vars.apt_mirror_nthreads / 1000 }) }}</div>
    </FieldNumber>
  </div>
</template>

<script>
function trimSlash(str) {
  if (str[str.length - 1] === '/') {
    return str.slice(0, str.length - 1)
  }
  return str
}

export default {
  props: {
    os_mirror: { type: Object, required: true },
  },
  data() {
    return {
      releases: {
        impish: '21.10',
        hirsute: '21.04',
        focal: '20.04',
        bionic: '18.04',
        xenial: '16.04',
        trusty: '14.04',
      },
      architecture: {
        ubuntu: {
          amd64: 'amd64',
          i386: 'i386',
          src: 'src',
        },
        'docker_ubuntu': {
          amd64: 'amd64',
          arm64: 'arm64',
          armhf: 'armhf',
          ppc64el: 'ppc64el',
          s390x: 's390x',
        }
      },
      mirrorOptions: {
        ubuntu: [
          'https://repo.huaweicloud.com/ubuntu',
          'https://mirrors.aliyun.com/ubuntu',
          'https://mirrors.tuna.tsinghua.edu.cn/ubuntu',
          'http://cn.archive.ubuntu.com/ubuntu',
          'https://mirrors.cloud.tencent.com/ubuntu',
          'http://ftp.sjtu.edu.cn/ubuntu',
          'http://mirrors.163.com/ubuntu',
          'http://mirrors.nju.edu.cn/ubuntu',
        ],
        'docker_ubuntu': [
          'https://repo.huaweicloud.com/docker-ce/linux/ubuntu',
          'https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu',
          'https://mirrors.aliyun.com/docker-ce/linux/ubuntu',
          'https://mirrors.cloud.tencent.com/docker-ce/linux/ubuntu',
        ]
      }
    }
  },
  inject: ['editMode'],
  computed: {
    vars: {
      get () { return this.os_mirror.inventory.all.children.target.vars },
      set () {},
    },
    prop () {
      return 'inventory.all.children.target.vars'
    },
    apt_mirror_ubuntu_mirror: {
      get () {
        if (this.vars.apt_mirror_ubuntu_mirror_protocol) {
          return this.vars.apt_mirror_ubuntu_mirror_protocol + this.vars.apt_mirror_ubuntu_mirror
        }
        return ''
      },
      set (v) {
        if (v) {
          let splited = v.split("://")
          this.vars.apt_mirror_ubuntu_mirror_protocol = splited[0] + "://"
          this.vars.apt_mirror_ubuntu_mirror = trimSlash(splited[1])
        } else {
          this.vars.apt_mirror_ubuntu_mirror_protocol = undefined
          this.vars.apt_mirror_ubuntu_mirror = undefined
        }
      }
    },
    apt_mirror_architecture: {
      get () {
        let temp = {}
        for (let i in this.vars.apt_mirror_repos) {
          let repo = this.vars.apt_mirror_repos[i]
          for (let key in this.architecture[this.os_mirror.status.type]) {
            if (repo.indexOf('deb-' + key) === 0) {
              temp[key] = true
            }
          }
        }
        let result = []
        for (let key in temp) {
          result.push(key)
        }
        return result
      },
      set (v) {
        if (v.length === 0) {
          this.$message.error('必须至少选一项')
          return
        }
        this.update_apt_mirror_repos(this.apt_mirror_repos, v)
      }
    },
    apt_mirror_repos: {
      get () {
        let temp = {}
        for (let i in this.vars.apt_mirror_repos) {
          let repo = this.vars.apt_mirror_repos[i]
          for (let key in this.releases) {
            if (repo.indexOf(key) > 0) {
              temp[key] = true
            }
          }
        }
        let result = []
        for (let key in temp) {
          result.push(key)
        }
        return result
      },
      set (v) {
        if (v.length === 0) {
          this.$message.error('必须至少选一项')
          return
        }
        this.update_apt_mirror_repos(v, this.apt_mirror_architecture)
      }
    }
  },
  watch: {
    'os_mirror': function() {
      this.initDefaults()
    },
  },
  components: { },
  mounted () {
    this.initDefaults()
  },
  methods: {
    update_apt_mirror_repos(releases, architecture) {
      let temp = []
      let releases_not_empty = releases
      if (releases_not_empty.length === 0) {
        releases_not_empty = ['focal']
      }
      let architecture_not_empty = architecture
      if (architecture_not_empty.length === 0) {
        architecture_not_empty = ['amd64']
      }

      // FIXME 需要确保 gpg 文件被正确下载
      for (let r of releases_not_empty) {
        for (let a of architecture_not_empty) {
          if (this.os_mirror.status.type === 'ubuntu') {
            temp.push(`deb-${a} {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${r} main restricted universe multiverse`)
            temp.push(`deb-${a} {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${r}-backports main restricted universe multiverse`)
            temp.push(`deb-${a} {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${r}-security main restricted universe multiverse`)
            temp.push(`deb-${a} {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${r}-updates main restricted universe multiverse`)
            temp.push(`deb-${a} {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${r}-proposed main restricted universe multiverse`)
          } else if (this.os_mirror.status.type === 'docker-ce_ubuntu') {
            temp.push(`deb-${a} {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${r} stable`)
          }
        }
      }
      this.vars.apt_mirror_repos = temp
    },
    initDefaults () {
      let vars = this.vars
      vars.apt_mirror_schedule = vars.apt_mirror_schedule || [ 
        {
          name: 'apt-mirror updates',
          special_time: 'daily',
          job: '/usr/bin/apt-mirror > {{ apt_mirror_dir }}/var/cron.log',
          cron_file: 'apt_mirror',
        } 
      ]
      vars.apt_mirror_repos = vars.apt_mirror_repos || []
      if (vars.apt_mirror_populate_repos === undefined) {
        vars.apt_mirror_populate_repos = true
      }
      if (vars.apt_mirror_enable_limit_rate === undefined) {
        vars.apt_mirror_enable_limit_rate = false
      }
      vars.apt_mirror_limit_rate = vars.apt_mirror_limit_rate || 125
      vars.apt_mirror_nthreads = vars.apt_mirror_nthreads || 10
      // vars.apache2_default_port = vars.apache2_default_port || 80
      vars.apt_mirror_repos = vars.apt_mirror_repos || [
        "deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-backports main restricted universe multiverse",
        "deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-security main restricted universe multiverse",
        "deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-updates main restricted universe multiverse",
        "deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} focal-proposed main restricted universe multiverse",
      ]
    },
    async loadSpecialTimeOptions () {
      return [
        { label: this.t('daily'), value: 'daily' },
        { label: this.t('weekly'), value: 'weekly' },
        { label: this.t('monthly'), value: 'monthly' },
      ]
    }
  }
}
</script>

<style scoped lang="css">
.placeholder {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
