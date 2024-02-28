<template>
  <span class="avatar" :style="{ width: width + 'px', height: width + 'px' }">
<!--    <img-->
<!--      :src="-->
<!--        avatar && avatar != ''-->
<!--          ? avatar-->
<!--          : `${proxy.globalInfo.avatarUrl}${userId}?${timestamp}`-->
<!--      "-->
<!--      v-if="userId"-->
<!--    />-->
    <img :src="IMG" v-if="userId" alt="头像">
  </span>
</template>

<script setup>
import { getCurrentInstance } from "vue";
import IMG from '../assets/qq.png'
const { proxy } = getCurrentInstance();

const props = defineProps({
  userId: {
    type: String,
  },
  avatar: {
    type: String,
  },
  timestamp: {
    type: Number,
    default: 0,
  },
  width: {
    type: Number,
    default: 40,
  },
});
</script>

<style lang="scss" scoped>
.avatar {
  display: flex;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  img {
    width: 100%;
    object-fit: cover;
  }
}
</style>