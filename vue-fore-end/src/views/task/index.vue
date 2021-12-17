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
          <el-tag v-for="item in scope.row.scan_area" :key="item" class="tag">{{
            item
          }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column
        header-align="center"
        align="center"
        label="扫描状态"
      >
      <template slot-scope="scope">
        <el-tag
          v-if="scope.row.status == '运行中'"
          type="success"
          size="small"
          style="margin-right: 10px;"
        >
          正在扫描
        </el-tag>
        <el-tag
          v-else-if="scope.row.status == 2"
          type="success"
          size="small"
          style="margin-right: 10px;"
        >
          扫描完成
        </el-tag>
        <el-tag
          v-else-if='scope.row.status == "停止"'
          type="warning"
          size="small"
          style="margin-right: 10px;"
        >
          手动停止
        </el-tag>
      </template>
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
      <el-table-column prop="address" label="操作" width="300">
        <template slot-scope="scope">
          <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button
            size="mini"
            type="primary"
            @click="handleStart(scope.row.id)"
            :disabled = "scope.row.status == '运行中'"
            >开始</el-button
          >
          <el-button size="mini" type="warning"
          @click="handleStop(scope.row.id)"
          :disabled = "scope.row.status == '停止'"
          >停止</el-button>
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
      <TaskForm :form="task_form" :company_list="company_list"></TaskForm>
      <div slot="footer" class="dialog-footer">
        <el-button @click="is_show_create_dialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateTask">确认</el-button>
      </div>
    </el-dialog>

    <!-- 修改任务信息的对话框 -->
    <el-dialog title="修改任务" :visible.sync="is_show_update_dialog">
      <TaskForm :form="task_form" :company_list="company_list"></TaskForm>
      <div slot="footer" class="dialog-footer">
        <el-button @click="is_show_update_dialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateTask">确认</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  addTask,
  getTaskList,
  deleteTask,
  startTask,
  stopTask,
  editTask,
} from "@/api/task";
import { getCompanyList } from "@/api/company";
import { Message } from "element-ui";
import TaskForm from "@/components/TaskForm";
export default {
  components: {
    TaskForm,
  },
  data() {
    return {
      company_list: [],
      task_list: [],
      page: {
        offset: 0,
        count: 10,
      },
      is_show_update_dialog: false,
      is_show_create_dialog: false,
      task_form: {
        name: "",
        scan_area: [],
        company_id: "",
      },
    };
  },
  created() {
    this.getCompanyList();
    getTaskList(this.page)
      .then((resp) => {
        this.task_list = resp.data.data;
      })
      .catch((err) => {
        console.log(err);
      });
  },
  methods: {
    getCompanyList() {
      getCompanyList(0, -1)
        .then((resp) => {
          //修改前端维护的数据
          this.company_list = resp.data.data;
        })
        .catch((err) => {
          console.log(err);
        });
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
          this.task_list.forEach((item) => {
            if (item.id == id) {
              item.status = "运行中";
            }
          });
        })
        .catch((err) => {
          console.log(err);
        });
    },
    //停止任务
    handleStop(id) {
      stopTask(id)
        .then((resp) => {
          Message({
            type: "success",
            message: "停止任务成功",
          });
          //TODO 修改前端维护的数据
          this.task_list.forEach((item) => {
            if (item.id == id) {
              item.status = "停止";
            }
          });
        })
        .catch((err) => {
          console.log(err);
        });
    },
    //打开创建任务对话框
    handleOpenCreateTaskDialog() {
      this.is_show_create_dialog = true;
    },
    //打开修改任务对话框
    handleEdit(row) {
      this.task_form = row;
      this.is_show_update_dialog = true;
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
    //添加新任务
    handleCreateTask() {
      addTask(this.task_form)
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
    //修改任务
    handleUpdateTask() {
      editTask(this.task_form)
        .then((resp) => {
          //更新前端维护的数据
          this.task_list = this.task_list.map((item) => {
            if (item.id === resp.data.data.id) {
              return resp.data.data;
            } else {
              return item;
            }
          });
          //关闭对话框
          this.is_show_update_dialog = false;
          //提示消息
          Message({
            type: "success",
            message: "修改成功",
          });
        })
        .catch((err) => {
          console.log(err);
        });
    },
  },
};
</script>
