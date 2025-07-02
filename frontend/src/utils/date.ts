/**
 * 通用时间格式化函数
 * @param dateStr 日期字符串或Date对象
 * @returns 格式化后的时间字符串
 */
export function formatDate(dateStr: string | Date): string {
  if (!dateStr) return '-';
  try {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return '-';

    const now = new Date();

    // 判断是否同一天（自然日）
    const isSameDay =
      date.getFullYear() === now.getFullYear() &&
      date.getMonth() === now.getMonth() &&
      date.getDate() === now.getDate();

    const oneDay = 24 * 60 * 60 * 1000;
    // 计算自然日差，使用拷贝，避免修改原对象
    const dateZero = new Date(date);
    dateZero.setHours(0, 0, 0, 0);
    const nowZero = new Date(now);
    nowZero.setHours(0, 0, 0, 0);
    const daysDiff = Math.floor((nowZero.getTime() - dateZero.getTime()) / oneDay);

    if (isSameDay) {
      // 今天内
      return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
    } else if (daysDiff < 7) {
      // 一周内
      return `${daysDiff}天前`;
    } else {
      // 超过一周
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  } catch {
    return '-';
  }
}

/**
 * 完整时间格式化函数（包含年月日时分秒）
 * @param dateStr 日期字符串或Date对象
 * @returns 完整格式的时间字符串
 */
export function formatFullDate(dateStr: string | Date): string {
  if (!dateStr) return '-';
  
  try {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) {
      return '-';
    }
    
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  } catch (error) {
    console.error('Full date formatting error:', error);
    return '-';
  }
}

/**
 * 仅日期格式化函数（年月日）
 * @param dateStr 日期字符串或Date对象
 * @returns 仅日期的字符串
 */
export function formatDateOnly(dateStr: string | Date): string {
  if (!dateStr) return '-';
  
  try {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) {
      return '-';
    }
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });
  } catch (error) {
    console.error('Date only formatting error:', error);
    return '-';
  }
} 
