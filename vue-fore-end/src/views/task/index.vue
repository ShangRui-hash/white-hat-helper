<template>
  <div class="app-container">
    <el-button type="primary" size="small" @click="handleOpenCreateTaskDialog"
      >添加任务</el-button
    >
    <el-table :data="task_list" style="width: 100%">
      <el-table-column prop="id" type="index" label="#" width="50">
      </el-table-column>
      <el-table-column prop="name" label="名称" width="100"> </el-table-column>
      <el-table-column
        header-align="center"
        align="left"
        label="扫描范围"
        width="400"
      >
      <template slot-scope="scope">
        <el-tag v-for="item in scope.row.scan_area" :key="item" class="tag"  >{{item}}</el-tag>
      </template>
      </el-table-column>
      <el-table-column
        header-align="center"
        align="center"
        prop="status"
        label="扫描状态"
      >
      </el-table-column>
      <el-table-column
        header-align="center"
        align="center"
        prop="company_name"
        label="公司名"
      >
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="150">
      </el-table-column>
      <el-table-column prop="address" label="操作" width="250">
        <template slot-scope="scope">
          <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button
            size="mini"
            type="primary"
            @click="handleStart(scope.row.id)"
            >开始</el-button
          >
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
    <!-- 创建项目的对话框 -->
    <el-dialog title="添加任务" :visible.sync="is_show_create_dialog">
      <el-form
        :model="create_task_form"
        :rules="create_task_rules"
        ref="create_task_form"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input
            v-model="create_task_form.name"
            placeholder="请输入任务名称"
          ></el-input>
        </el-form-item>
        <el-form-item label="所属公司" prop="company_name">
          <el-select
            v-model="create_task_form.company_id"
            placeholder="请选择所属公司"
            style="width: 100%"
          >
            <el-option
              v-for="item in company_list"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <FormTags
          @update_list="updateScanArea"
          label="扫描范围"
          closable="false"
          :list="create_task_form.scan_area"
          formLabelWidth="80px"
        ></FormTags>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="is_show_create_dialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateTask">确认</el-button>
      </div>
    </el-dialog>

    <!-- 修改公司信息的对话框 -->
    <el-dialog title="修改公司信息" :visible.sync="is_show_update_dialog">
      <el-form label-width="80px">
        <el-form-item label="名称">
          <el-input></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="is_show_update_dialog = false">取消</el-button>
        <el-button type="primary" @click="updateCompany">确认</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { addTask, getTaskList, deleteTask, startTask } from "@/api/task";
import { getCompanyList } from "@/api/company";
import { Message } from "element-ui";
import  FormTags  from "@/components/FormTags";
export default {
  components: {
     FormTags
 },
  data() {
    return {
      task_list: [
        {
          id: "1",
          name: "扫描任务1",
          scan_area: "192.168.2.20,192.168.2.21",
          status: "扫描中",
          company_name: "联想",
          created_at: "2021-12-1",
        },
      ],
      page: {
        offset: 0,
        count: 10,
      },
      is_show_update_dialog: false,
      is_show_create_dialog: false,
      create_task_form: {
        name: "",
        scan_area: [],
        company_id: "",
      },
      company_list: [
        {
          id: "1",
          name: "联想",
        },
        {
          id: "2",
          name: "华为",
        },
      ],
    };
  },
  created() {
    getCompanyList(0, -1)
      .then((resp) => {
        //修改前端维护的数据
        this.company_list = resp.data.data;
      })
      .catch((err) => {
        console.log(err);
      });
    getTaskList(this.page)
      .then((resp) => {
        this.task_list = resp.data.data;
      })
      .catch((err) => {
        console.log(err);
      });
  },
  methods: {
    updateScanArea(list) {
      this.create_task_form.scan_area = list;
    },
    //开始任务
    handleStart(id) {
      startTask(id)
        .then((resp) => {
          Message({
            type: "success",
            message: "开始任务成功",
          });
          //TODO 修改前端维护的数据
        })
        .catch((err) => {
          console.log(err);
        });
    },
    //打开创建任务对话框
    handleOpenCreateTaskDialog() {
      this.is_show_create_dialog = true;
    },
    //添加新任务
    handleCreateTask() {
      addTask(this.create_task_form)
        .then((resp) => {
          //更新前端维护的数据
          this.task_list.push(resp.data.data);
          //关闭对话框
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
    handleEdit(row) {
      console.log(row);
    },
    handleDelete(id) {
      deleteTask(id).then((resp) => {
        //更新前端维护的数据
        this.task_list = this.task_list.filter((item) => item.id !== id);
        //提示消息
        Message({
          type: "success",
          message: "删除成功",
        });
      });
    },
    updateCompany() {
      this.is_show_update_dialog = false;
    },
  },
};
</script>
