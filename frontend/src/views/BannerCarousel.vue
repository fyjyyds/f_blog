<template>
  <div class="banner-container">
    <div class="banner-carousel">
      <div
        v-for="(banner, idx) in banners"
        :key="idx"
        class="banner-slide"
        :class="{ active: currentIndex === idx }"
        :style="{ backgroundImage: `url(${getFullUrl(banner)})` }"
      >
        <div class="banner-overlay">
          <div class="banner-content">
            <h2 class="banner-title">欢迎来到 F.Blog</h2>
            <p class="banner-subtitle">分享你的想法，连接世界</p>
          </div>
        </div>
      </div>
    </div>
    <div class="banner-indicators" v-if="banners.length > 1">
      <div
        v-for="(banner, idx) in banners"
        :key="idx"
        class="indicator"
        :class="{ active: currentIndex === idx }"
        @click="setCurrentIndex(idx)"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { getBanners } from '../api/banner';

const banners = ref<any[]>([]);
const currentIndex = ref(0);
const backendBaseUrl = 'http://localhost:8080';

function getFullUrl(banner: any) {
  const path = banner?.image || banner;
  if (!path) return '';
  if (path.startsWith('http')) return path;
  return backendBaseUrl + path;
}

function setCurrentIndex(index: number) {
  currentIndex.value = index;
}

let timer: number;

onMounted(async () => {
  try {
    const res = await getBanners();
    banners.value = res.data.banners || res.data;
  } catch (error) {
    console.error('获取轮播图失败:', error);
  }
  timer = window.setInterval(() => {
    if (banners.value.length > 1) {
      currentIndex.value = (currentIndex.value + 1) % banners.value.length;
    }
  }, 5000);
});

onUnmounted(() => {
  clearInterval(timer);
});
</script>

<style scoped>
.banner-container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.banner-carousel {
  width: 100%;
  height: 280px;
  position: relative;
  overflow: hidden;
  border-radius: 12px;
}

.banner-slide {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  opacity: 0;
  transition: opacity 0.6s ease;
}

.banner-slide.active {
  opacity: 1;
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(0,0,0,0.35), rgba(0,0,0,0.15));
  display: flex;
  align-items: center;
  justify-content: center;
}

.banner-content {
  text-align: center;
  color: #fff;
}

.banner-title {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 8px rgba(0,0,0,0.3);
}

.banner-subtitle {
  font-size: 16px;
  margin: 0;
  opacity: 0.9;
  text-shadow: 0 1px 4px rgba(0,0,0,0.3);
}

.banner-indicators {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 12px;
}

.indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #d1d5db;
  cursor: pointer;
  transition: all 0.2s;
}

.indicator.active {
  background: #5b6abf;
  width: 20px;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .banner-carousel {
    height: 200px;
  }
  .banner-title {
    font-size: 22px;
  }
}
</style>
