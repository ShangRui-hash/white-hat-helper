<template>
  <div class="app-container">
    <el-skeleton v-show="is_show_skeleton" :rows="18" animated />
    <el-card v-for="host in hosts" :key="host.id" class="box-card">
      <div class="w-100percent">
        <!-- 操作系统 -->
        <el-tag v-if="host.os.length > 0" type="info" class="tag">{{
          host.os
        }}</el-tag>
        <!-- ip -->
        <el-link
          type="primary"
          style="font-size: 18px"
          @click="handleOpenDetail(host.ip)"
          >{{ host.ip }}</el-link
        >
        <!-- 域名 -->
        <div>
          <el-tag
            class="tag"
            type="primary"
            v-for="domain in host.domain_list"
            :key="domain"
            >{{ domain }}</el-tag
          >
        </div>

        
        <!-- 端口 -->
        <div
          v-if="host.ports && host.ports.length > 0"
          class="flex m-10 w-100percent"
        >
          <div
            type="primary"
            v-for="port in host.ports.slice(0, 20)"
            :key="port"
            class="service"
          >
            <a href="">{{ port.port }}/{{ port.service }}</a>
          </div>
          <div v-show="host.ports && host.ports.length > 20">......</div>
        </div>
        <!-- web服务 -->
        <el-table
          v-if="host.webs && host.webs.length > 0"
          :data="host.webs.slice(0, 5)"
          style="width: 100%"
        >
          <el-table-column prop="url" label="url"> </el-table-column>
          <el-table-column prop="status_code" label="status_code">
          </el-table-column>
          <el-table-column prop="title" label="title"> </el-table-column>
          <el-table-column prop="location" label="location"></el-table-column>
          <el-table-column prop="fingerprint" label="fingerprint">
            <template slot-scope="scope" >
              <el-tag
                v-for="item in scope.row.fingerprint"
                :key="item"
                class="tag"
                >{{ item }}</el-tag
              >
            </template>
          </el-table-column>
        </el-table>
        <el-link class="m-l-10" v-show="host.webs && host.webs.length > 10"
          >.......</el-link
        >
      </div>
    </el-card>
    <load-more-btn
      :loading="loading"
      :loadmore_btn_text="loadmore_btn_text"
      @loadmore="loadmore"
    ></load-more-btn>
  </div>
</template>

<script>
import LoadMoreBtn from "@/components/LoadMoreBtn";
import { getHostList } from "@/api/hosts";
export default {
  name: "HostsList",
  components: {
    LoadMoreBtn,
  },
  data() {
    return {
      is_show_skeleton: false,
      loadmore_btn_text: "",
      loading: false,
      hosts: [],
      page: {
        offset: 0,
        count: 10,
      },
      search: {
        name: "",
        ip: "",
        status: "",
      },
      sort: {
        sort_by: "name",
        sort_type: "asc",
      },
    };
  },
  created() {
    console.log("HostsList created");
    this.is_show_skeleton = true;
    this.loading = true;
    this.getHosts();
  },
  methods: {
    handleOpenDetail(ip) {
      this.$router.push({
        path: `/host/detail/${ip}`,
      });
    },
    //加载更多
    loadmore() {
      this.loading = true;
      this.getHosts();
    },
    //获取主机列表
    getHosts() {
      getHostList({
        company_id: this.$route.params.company_id,
        offset: this.page.offset,
        count: this.page.count,
      })
        .then((resp) => {
          console.log(resp);
          let hosts = resp.data.data == null ? [] : resp.data.data;
          //1.遍历获取到的所有主机
          for (let i = 0; i < hosts.length; i++) {
            //2.寻找现有主机中有无该主机
            let index = this.hosts.findIndex((item) => item.ip == hosts[i].ip);
            if (index == -1) {
              //3.如果没有，则添加
              this.hosts.push(hosts[i]);
            } else {
              //4.如果有，则更新
              this.hosts[index] = hosts[i];
            }
          }
          //5.更新分页信息
          this.page.offset = this.hosts.length;
          this.loading = false;
          this.is_show_skeleton = false;
          //6.更新按钮文字
          if (hosts.length < this.count) {
            this.loadmore_btn_text = this.config.nomore_text;
          } else {
            this.loadmore_btn_text = this.config.loadmore_text;
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
  },
};
</script>
