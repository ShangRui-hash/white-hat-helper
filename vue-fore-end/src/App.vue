<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: "App",
};
</script>

<style>
.tag {
  margin-right: 10px;
  margin-top: 10px;
}
.flex {
  display: flex;
  flex-wrap: wrap;
}
.flex-space-between {
  justify-content: space-between;
}
.service {
  flex: 0 1 auto;
  min-width: 20%;
}
.service a,a:hover,a:focus {
  color: #409eff;
  font-size: 14px;
}
.service a:hover {
  text-decoration:underline;
}
.m-20 {
  margin:20px;
}
.m-10{
  margin:10px;
}
.w-100percent{
  width: 100%;
}
.m-l-10{
  margin-left:10px;
}
.m-t-15{
  margin-top:15px;
}
.m-tb-15{
  margin-top:15px;
  margin-bottom:15px;
}
.code{
  width:100%;
}
.border{
  border:1px solid #E4E7ED;
  padding:15px;
}
</style>
