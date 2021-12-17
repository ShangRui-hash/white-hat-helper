<template>
  <el-form
    :model="form"
    :rules="rules"
    ref="create_task_form"
    label-width="80px"
  >
    <el-form-item label="名称" prop="name">
      <el-input v-model="form.name" placeholder="请输入任务名称"></el-input>
    </el-form-item>
    <el-form-item label="所属公司" prop="company_name">
      <el-select
        v-model="form.company_id"
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
      :list="form.scan_area"
      formLabelWidth="80px"
    ></FormTags>
  </el-form>
</template>

<script>
import FormTags from "@/components/FormTags";
export default {
  props: ["form","company_list"],
  components: {
    FormTags,
  },
  data() {
    return {
      company_list: [],
      form: {
        name: "",
        scan_area: [],
        company_id: "",
      },
    };
  },
  methods: {
    updateScanArea(list) {
      this.task_form.scan_area = list;
    },
  },
};
</script>
