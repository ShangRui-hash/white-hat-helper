<template>
  <div class="app-container documentation-container">
    <el-button
      type="primary"
      size="small"
      @click="handleOpenCreateCompanyDialog"
      >添加公司</el-button
    >
    <el-table :data="company_list" style="width: 100%">
      <el-table-column prop="id" type="index" label="#" width="50">
      </el-table-column>
      <el-table-column prop="name" label="名称" width="100">
        <template slot-scope="scope">
          <el-link
            type="primary"
            @click="handleClickCompanyName(scope.row.id)"
            >{{ scope.row.name }}</el-link
          >
        </template>
      </el-table-column>
      <el-table-column
        header-align="center"
        align="center"
        prop="asset_count"
        label="资产数"
      >
      </el-table-column>
      <el-table-column
        header-align="center"
        align="center"
        prop="task_count"
        label="任务数"
      >
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="150">
      </el-table-column>
      <el-table-column prop="address" label="操作">
        <template slot-scope="scope">
          <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
          <span style="margin-right: 10px"></span>
          <el-popconfirm
            title="您确定要删除吗?"
            @confirm="handleDelete(scope.row.id)"
          >
            <el-button slot="reference" size="mini" type="danger"
              >删除</el-button
            >
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <load-more-btn
      :loading="loading"
      :loadmore_btn_text="loadmore_btn_text"
      @loadmore="loadmore"
    ></load-more-btn>
    <!-- 添加公司对话框 -->
    <el-dialog
      title="添加公司"
      :visible.sync="is_show_create_dialog"
      width="30%"
    >
      <el-form ref="form" :model="create_company_form" label-width="80px">
        <el-form-item label="公司名称" prop="name">
          <el-input
            v-model="create_company_form.name"
            placeholder="公司名称"
          ></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="is_show_create_dialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateCompany">确认</el-button>
      </div>
    </el-dialog>
    <!-- 修改公司信息的对话框 -->
    <el-dialog title="修改公司信息" :visible.sync="is_show_update_dialog">
      <el-form label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="update_company_form.name"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="is_show_update_dialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateCompany">确认</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  addCompany,
  deleteCompany,
  getCompanyList,
  updateCompany,
} from "@/api/company";
import LoadMoreBtn from "@/components/LoadMoreBtn";
import { Message } from "element-ui";
export default {
  components: {
    LoadMoreBtn,
  },
  data() {
    return {
      loadmore_btn_text: "",
      loading: false,
      offset: 0,
      count: 10,
      is_show_update_dialog: false,
      is_show_create_dialog: false,
      create_company_form: {
        name: "",
      },
      update_company_form: {
        id: "",
        name: "",
      },
      company_list: [],
    };
  },
  created() {
    this.loading = true;
    //获取公司列表
    this.getCompanyList();
  },
  methods: {
    loadmore() {
      this.loading = true;
      this.getCompanyList();
    },
    //获取公司列表
    getCompanyList() {
      getCompanyList(this.offset, this.count)
        .then((res) => {
          //更新前端维护的数据
          let companies =  res.data.data==null?[]:res.data.data;
          for(let i = 0; i < companies.length; i++){
            let index = this.company_list.findIndex(item => item.id == companies[i].id);
            if (index == -1) {
              this.company_list.push(companies[i]);
            }else{
              this.company_list[index] = companies[i];
            }
          }
          //更新分页信息
          this.offset = this.company_list.length;
          this.loading = false;
          //更新按钮文字
          if (companies.length < this.count) {
            this.loadmore_btn_text = this.config.nomore_text;
          } else {
            this.loadmore_btn_text = this.config.loadmore_text;
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    //跳转到公司资产列表页
    handleClickCompanyName(id) {
      this.$router.push({
        path: `/host/list/${id}`,
      });
    },
    //打开添加公司对话框
    handleOpenCreateCompanyDialog() {
      this.is_show_create_dialog = true;
    },
    //添加公司
    handleCreateCompany() {
      addCompany(this.create_company_form)
        .then((res) => {
          //更新前端维护的数据
          this.company_list.push(res.data.data);
          //关闭模态框
          this.is_show_create_dialog = false;
          //提示消息
          Message({
            type: "success",
            message: "添加成功",
          });
        })
        .catch((err) => {
          console.log(err);
        });
    },
    //点击编辑按钮
    handleEdit(row) {
      //1.显示模态框
      this.is_show_update_dialog = true;
      //2.设置模态框的数据
      this.update_company_form.name = row.name;
      this.update_company_form.id = row.id;
    },
    //确认修改
    handleUpdateCompany() {
      //1.前端验证
      if (!this.update_company_form.name) {
        Message({
          type: "error",
          message: "请输入公司名称",
          duration: 5 * 1000,
        });
        return;
      }
      //2.发送请求
      updateCompany(this.update_company_form)
        .then((res) => {
          //更新前端维护的数据
          let index = this.company_list.findIndex((item) => {
            return item.id === this.update_company_form.id;
          });
          this.company_list[index].name = this.update_company_form.name;
          //关闭模态框
          this.is_show_update_dialog = false;
          //提示消息
          Message({
            type: "success",
            message: "修改成功",
            duration: 5 * 1000,
          });
        })
        .catch((err) => {
          console.log(err);
        });
    },
    //删除
    handleDelete(id) {
      deleteCompany(id).then(
        (res) => {
          //更新前端维护的数据
          this.company_list = this.company_list.filter(
            (item) => item.id !== id
          );
          //提示消息
          Message({
            type: "success",
            message: "删除成功",
            duration: 5 * 1000,
          });
        },
        (err) => {
          console.log(err);
        }
      );
    },
  },
};
</script>

<style lang="scss" scoped></style>
