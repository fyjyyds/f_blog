<template>
  <div class="banner-container">
    <div class="banner-carousel">
      <div 
        class="banner-slide"
        v-for="(banner, idx) in banners" 
        :key="idx"
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
      
      <div class="banner-indicators">
        <button 
          v-for="(banner, idx) in banners" 
          :key="idx"
          class="indicator"
          :class="{ active: currentIndex === idx }"
          @click="setCurrentIndex(idx)"
        ></button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getBanners } from '../api/banner';

const banners = ref<string[]>([]);
const currentIndex = ref(0);
const backendBaseUrl = 'http://localhost:8080';

function getFullUrl(path: string) {
  if (path.startsWith('http')) return path;
  return backendBaseUrl + path;
}

function setCurrentIndex(index: number) {
  currentIndex.value = index;
}

onMounted(async () => {
  try {
    const res = await getBanners();
    banners.value = res.data.banners || res.data;
  } catch (error) {
    console.error('获取轮播图失败:', error);
  }
});
</script>

<style scoped>
.banner-container {
  width: 100%;
  height: 300px;
  position: relative;
  overflow: hidden;
  margin-bottom: 32px;
}

.banner-carousel {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.banner-slide {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  opacity: 0;
  transition: opacity 0.8s ease-in-out;
  transform: scale(1.1);
}

.banner-slide.active {
  opacity: 1;
  transform: scale(1);
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    135deg,
    rgba(0, 0, 0, 0.4) 0%,
    rgba(0, 0, 0, 0.2) 50%,
    rgba(0, 0, 0, 0.6) 100%
  );
  display: flex;
  align-items: center;
  justify-content: center;
}

.banner-content {
  text-align: center;
  color: white;
  padding: 40px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.banner-title {
  font-size: 36px;
  font-weight: 800;
  margin: 0 0 12px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.banner-subtitle {
  font-size: 18px;
  margin: 0;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 400;
}

.banner-indicators {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 12px;
  z-index: 10;
}

.indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.5);
  background: transparent;
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.8);
}

.indicator.active {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-color: #667eea;
  box-shadow: 0 0 10px rgba(102, 126, 234, 0.5);
}

@media (max-width: 768px) {
  .banner-container {
    height: 200px;
    margin-bottom: 16px;
  }
  
  .banner-content {
    padding: 20px;
  }
  
  .banner-title {
    font-size: 24px;
  }
  
  .banner-subtitle {
    font-size: 14px;
  }
}
</style>