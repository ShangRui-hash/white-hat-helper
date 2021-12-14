<template>
  <div class="app-container">
    <el-skeleton v-show="is_show_skeleton" :rows="18" animated />
    <div v-show="is_show_skeleton == false">
    <el-row :gutter="10">
      <el-card class="box-card">
        <el-descriptions title="基本信息" border>
          <el-descriptions-item label="IP">{{
            target.ip
          }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{
            target.os
          }}</el-descriptions-item>
          <el-descriptions-item label="所属公司">{{
            target.company
          }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </el-row>
    <el-row :gutter="10">
      <el-card class="box-card">
        <el-table :data="target.ports" style="width: 100%">
          <el-table-column prop="port" label="端口"> </el-table-column>
          <el-table-column prop="status" label="状态"> </el-table-column>
          <el-table-column prop="service" label="服务"> </el-table-column>
          <el-table-column prop="protocol" label="协议"> </el-table-column>
          <el-table-column prop="version" label="版本"> </el-table-column>
        </el-table>
      </el-card>
    </el-row>
    <el-card v-for="web in target.webs" :key="web.url" class="box-card">
      <div>
        <!-- 响应码 -->
        <el-tag type="success" class="tag">{{ web.status_code }}</el-tag>
        <!-- url -->
        <el-link type="primary" style="font-size: 18px" :href="web.url">{{
          web.url
        }}</el-link>
        <!-- 标题 -->
        <el-link style="margin-left: 10px" :href="web.url">{{
          web.title
        }}</el-link>
        <!-- 指纹信息+防火墙 -->
        <div style="margin-top: 10px">
          <el-tag v-for="item in web.fingerprint" :key="item" class="tag">{{
            item
          }}</el-tag>
          <el-tag v-if="web.waf_name.lenght > 0" type="danger" class="tag">{{
            web.waf_name
          }}</el-tag>
        </div>
        <!-- 响应头 -->

        <div class="m-t-15 border">
          <div v-for="(valueList, key) in web.resp_header" :key="key">
            <span style="color: #409eff; margin-right: 10px">{{ key }}:</span>
            <span v-for="value in valueList" :key="value">{{ value }}</span>
          </div>
        </div>
        <!-- 响应体 -->
        <div class="m-tb-15 border">
          <codemirror
            ref="myCm"
            :value="web.resp_body"
            :options="codeOption"
            class="code"
          ></codemirror>
        </div>
      </div>
      <el-descriptions border direction="vertical">
        <el-descriptions-item label="目录" :span="3">
          <el-table :data="web.dirs" style="width: 100%">
            <el-table-column prop="url" label="URL"> </el-table-column>
            <el-table-column prop="status_code" label="响应码">
            </el-table-column>
            <el-table-column prop="title" label="标题"> </el-table-column>
            <el-table-column prop="content_size" label="响应报文大小">
            </el-table-column>
          </el-table>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
    </div>

  </div>
</template>

<script>
import { getHostDetail } from "@/api/hosts";
import { codemirror } from "vue-codemirror";
require("codemirror/mode/python/python.js");
require("codemirror/addon/fold/foldcode.js");
require("codemirror/addon/fold/foldgutter.js");
require("codemirror/addon/fold/brace-fold.js");
require("codemirror/addon/fold/xml-fold.js");
require("codemirror/addon/fold/indent-fold.js");
require("codemirror/addon/fold/markdown-fold.js");
require("codemirror/addon/fold/comment-fold.js");
export default {
  components: {
    codemirror,
  },
  data() {
    return {
      is_show_skeleton: false,
      codeOption: {
        value: "",
        theme: "default", //主题
        indentUnit: 2,
        smartIndent: true,
        tabSize: 4,
        readOnly: true, //只读
        showCursorWhenSelecting: true,
        lineNumbers: false, //是否显示行数
        firstLineNumber: 1,
      },
      target: {},
    };
  },
  created() {
    this.is_show_skeleton=true
    this.getHostDetail();
  },
  methods: {
    getHostDetail() {
      getHostDetail({
        ip: this.$route.params.ip,
      })
        .then((res) => {
          this.target = res.data.data;
          this.is_show_skeleton=false
        })
        .catch((err) => {
          console.log(err);
        });
    },
  },
};
</script>

<style scoped></style>
