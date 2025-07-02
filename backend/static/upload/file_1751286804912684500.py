import numpy as np
import matplotlib.pyplot as plt

# 模拟数据
height = np.linspace(0, 15, 300)  # 高度范围 (km)
range_km = np.linspace(0, 100, 500)  # 距离范围 (km)
data = np.random.rand(len(height), len(range_km)) * 70  # 模拟的反射率 (dbz)

# 创建网格数据
X, Y = np.meshgrid(range_km, height)

# 绘制图像
plt.figure(figsize=(10, 6))
plt.pcolormesh(X, Y, data, cmap='jet', shading='auto', vmin=5, vmax=65)
plt.colorbar(label='dbz')  # 添加颜色条

# 添加标签和标题
plt.title('dbz')
plt.xlabel('Range (km)')
plt.ylabel('Height (km)')

# 设置范围
plt.ylim(0, 15)
plt.xlim(0, 100)

# 显示图像
plt.show()