<template>
  <div class="flex flex-center" style="padding:15px;">
    <el-button
      style="width:100%"
      type="primary"
      plain
      size="medium"
      @click="loadmore"
      :loading="loading"
      :disabled="loadmore_btn_text == config.nomore_text"
    >{{loadmore_btn_text}}</el-button>
  </div>
</template>

<script>
export default {
  props: ["loading", "loadmore_btn_text"],
  methods: {
    loadmore() {
      this.$emit("loadmore");
    }
  }
};
</script>