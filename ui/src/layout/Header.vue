<template>
  <div class="header">
    <h2 class="left">
      <a href="https://www.zendata.cn" target="_blank">{{$t('site.title')}}</a>
    </h2>
    <div class="center"></div>
    <div class="right">
      <span class="dir">{{$t('msg.workdir')}}: {{workDir}}</span>

      <select-lang :prefixCls="'select-lang'" />

      <a href="https://www.zendata.cn/book/zendata/" target="_blank">{{$t('msg.help')}}</a>
    </div>
  </div>
</template>

<script>

import {getWorkDir} from "../api/manage";
import {config} from "../utils/vari";
import SelectLang from '../components/SelectLang'

export default {
  name: 'Header',
  components: {
    SelectLang
  },
  data () {
    return {
      workDir: '',
    }
  },
  computed: {
  },
  created () {
    getWorkDir().then(json => {
      console.log('getWorkDir', json)
      const that = this
      that.defs = json.data
      this.workDir = json.workDir
      config.workDir = this.workDir
    })
  },
  mounted () {
  },
  methods: {
  }
}
</script>

<style lang="less" scoped>
.header {
  display: flex;
  height: 49px;
  line-height: 49px;

  a {
    color: #fff;
  }

  .left {
    margin: 0 15px;
    width: 100px;
    color: #fff;
  }
  .center {
    flex: 1;
  }
  .right {
    margin: 0 15px;
    width: 500px;
    text-align: right;

    .dir {
      display: inline-block;
      padding-right: 20px;
    }
  }

  .select-lang {
    padding-right: 10px;
    cursor: pointer;
  }
}
</style>
