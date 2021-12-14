<template>
  <el-form-item v-show='(list && list.length>0) || closable' :label="label" :label-width="formLabelWidth">
    <el-tag
      :key="item"
      v-for="item in list"
      :closable="closable"
      :disable-transitions="false"
      @close="deleteItem(item)"
    >{{item}}</el-tag>
    <span v-if="closable">
      <el-input
        class="input-new-tag"
        v-if="input_visible"
        v-model="input_value"
        ref="saveTagInput"
        size="small"
        @keyup.enter.native="add"
        @blur="add"
      ></el-input>
      <el-button v-else class="button-new-tag" size="small" @click="showInput">+</el-button>
    </span>
  </el-form-item>
</template>

<script>
export default {
  props: ["formLabelWidth", "label", "list", "closable"],
  data: () => {
    return {
      input_value: "",
      input_visible: false
    };
  },
  methods: {
    //从参数列表中移除一个参数
    deleteItem(item) {
      this.list.splice(this.list.indexOf(item), 1);
      //通知父组件删除元素
      this.$emit("update_list", this.list);
    },
    //将input框中的标签添加到标签列表中
    add() {
      let inputValue = this.input_value;
      if (inputValue) {
        this.list.push(inputValue);
      }
      this.input_visible = false;
      this.input_value = "";
      //通知父组件添加元素
      this.$emit("update_list", this.list);
    },
    //显示添加参数的输入框
    showInput() {
      this.input_visible = true;
      this.$nextTick(() => {
        this.$refs.saveTagInput.$refs.input.focus();
      });
    }
  }
};
</script>


<style scoped>
.el-tag + .el-tag {
  margin-left: 10px;
}
.button-new-tag {
  margin-left: 10px;
  height: 32px;
  line-height: 30px;
  padding-top: 0;
  padding-bottom: 0;
}
.input-new-tag {
  width: 90px;
  margin-left: 10px;
  vertical-align: bottom;
}
</style>