<template>
  <div class="graph-container" :class="{ 'filter-open': filterVisible }">
    <!-- トグルボタンを追加 -->
    <button class="toggle-filter-btn" @click="toggleFilterPanel">
      {{ filterVisible ? '◀' : '▶' }}
    </button>
    <!-- フィルターパネル（左側） - v-showで表示・非表示を切り替え -->
    <div class="filter-panel" :class="{ 'hidden': !filterVisible }">
      <h3>フィルター</h3>

      <!-- 表示期間ラジオボタン（チャートの上部に配置） -->
      <div class="filter-item">
        <label>集計軸</label>
        <div class="radio-group" v-if="filterVisible">
          <label class="custom-radio">
            <input type="radio" v-model="timeFilter" value="day" @change="updateDuration();updateChart()"/>
            <span class="radio-label">日</span>
          </label>
          <label class="custom-radio">
            <input type="radio" v-model="timeFilter" value="week" @change="updateDuration();updateChart()"/>
            <span class="radio-label">週</span>
          </label>
          <label class="custom-radio">
            <input type="radio" v-model="timeFilter" value="month" @change="updateDuration();updateChart()"/>
            <span class="radio-label">月</span>
          </label>
        </div>
        <!-- 最小化時のホバー表示用要素 -->
        <div class="hover-content time-filter-hover">
          <div class="radio-group">
            <label class="custom-radio">
              <input type="radio" v-model="timeFilter" value="day" @change="updateDuration();updateChart()"/>
              <span class="radio-label">日</span>
            </label>
            <label class="custom-radio">
              <input type="radio" v-model="timeFilter" value="week" @change="updateDuration();updateChart()"/>
              <span class="radio-label">週</span>
            </label>
            <label class="custom-radio">
              <input type="radio" v-model="timeFilter" value="month" @change="updateDuration();updateChart()"/>
              <span class="radio-label">月</span>
            </label>
          </div>
        </div>
      </div>

      <!-- 案件種別のチェックボックス（ラジオボタンから変更） -->
      <div class="filter-item">
        <label>案件種別</label>
        <div class="radio-group" v-if="filterVisible">
          <label class="custom-radio">
            <input
                type="checkbox"
                :checked="facilityTypes.includes(FacilityType.Ordered)"
                @change="toggleFacilityType(FacilityType.Ordered)"
                class="hidden-checkbox"
            />
            <span class="radio-label">{{ FacilityTypeMap[FacilityType.Ordered] }}</span>
          </label>
          <label class="custom-radio">
            <input
                type="checkbox"
                :checked="facilityTypes.includes(FacilityType.Prepared)"
                @change="toggleFacilityType(FacilityType.Prepared)"
                class="hidden-checkbox"
            />
            <span class="radio-label">{{ FacilityTypeMap[FacilityType.Prepared] }}</span>
          </label>
        </div>
        <!-- 最小化時のホバー表示用要素 -->
        <div class="hover-content facility-type-hover">
          <div class="radio-group">
            <label class="custom-radio">
              <input
                  type="checkbox"
                  :checked="facilityTypes.includes(FacilityType.Ordered)"
                  @change="toggleFacilityType(FacilityType.Ordered)"
                  class="hidden-checkbox"
              />
              <span class="radio-label">{{ FacilityTypeMap[FacilityType.Ordered] }}</span>
            </label>
            <label class="custom-radio">
              <input
                  type="checkbox"
                  :checked="facilityTypes.includes(FacilityType.Prepared)"
                  @change="toggleFacilityType(FacilityType.Prepared)"
                  class="hidden-checkbox"
              />
              <span class="radio-label">{{ FacilityTypeMap[FacilityType.Prepared] }}</span>
            </label>
          </div>
        </div>
      </div>

      <div class="filter-item">
        <label>期間選択</label>
        <select v-model="durationFilter" @change="updateChart" v-if="filterVisible">
          <template v-for="item in durationOptions" :key="item.key">
            <option :value="item.value">{{ item.key }}</option>
          </template>
        </select>
        <!-- 最小化時のホバー表示用要素 -->
        <div class="hover-content duration-filter-hover">
          <select v-model="durationFilter" @change="updateChart">
            <template v-for="item in durationOptions" :key="item.key">
              <option :value="item.value">{{ item.key }}</option>
            </template>
          </select>
        </div>
      </div>

      <div class="filter-item">
        <label>開始日</label>
        <input 
          type="date" 
          :value="formatDateForInput(startDate)" 
          @change="handleStartDateChange($event)"
          :min="formatDateForInput(globalStartTimestamp)"
          v-if="filterVisible"
        />
        <!-- 最小化時のホバー表示用要素 -->
        <div class="hover-content start-date-hover">
          <input 
            type="date" 
            :value="formatDateForInput(startDate)" 
            @change="handleStartDateChange($event)"
            :min="formatDateForInput(globalStartTimestamp)"
          />
        </div>
      </div>
    </div>
    <!-- グラフ表示エリア - フィルター非表示時に全幅表示 -->
    <div class="chart-area" :class="{ 'full-width': !filterVisible }">
      <div class="legend-buttons">
        <button class="legend-btn" @click="selectAllSeries">すべて選択</button>
        <button class="legend-btn" @click="unselectAllSeries">すべて解除</button>
        <a href="#" @click.prevent="updateChart" class="reload-btn">
          <span class="material-symbols-outlined">refresh</span>
          <span class="text">リロード</span>
        </a>
      </div>

      <apex-charts
          ref="chartRef"
          :options="chartOptions"
          :series="series"
          height="100%"
          @legendClick="handleLegendClick"
          style="flex-grow: 1;"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, reactive, ref, Ref} from 'vue'
import ApexCharts from 'vue3-apexcharts'
import {usePileUpGraph} from "@/composable/pileUpGraph";
import {FacilityType, FacilityTypeMap} from "@/const/common";
import dayjs from 'dayjs';


const {
  facilities,
  departments,
  pileUps,
  facilityTypes,
  timeFilter,
  series,
  xLabels,
  selectedSeries,
  refreshData,
  durationOptions,
  durationFilter,
  updateDuration,
  startDate,
  updateStartDate,
  globalStartTimestamp
} = await usePileUpGraph()
// 案件種別の切り替え関数
const toggleFacilityType = (type: FacilityType) => {
  const index = facilityTypes.indexOf(type);
  if (index === -1) {
    facilityTypes.push(type);
  } else {
    facilityTypes.splice(index, 1);
  }
  updateChart();
};

// チャートの参照を作成
const chartRef = ref<null | any>(null) as Ref<typeof ApexCharts>;

// フィルターパネルの表示状態を管理するリアクティブな変数
const filterVisible = ref(true);

// フィルターパネルの表示/非表示を切り替える関数
const toggleFilterPanel = () => {
  filterVisible.value = !filterVisible.value;
};
// すべてのシリーズの表示/非表示を切り替える関数
const toggleAllSeries = (visible: boolean) => {
  if (chartRef.value?.chart) {
    series.value.forEach(seriesItem => {
      if (visible) {
        chartRef.value.chart.showSeries(seriesItem.name);
      } else {
        chartRef.value.chart.hideSeries(seriesItem.name);
      }

      // 内部状態も更新
      const index = series.value.findIndex(s => s.name === seriesItem.name);
      if (index >= 0) {
        selectedSeries.value[index] = !visible;
      }
    });
  }
};

// 元の関数は以下のように置き換え
const selectAllSeries = () => toggleAllSeries(true);
const unselectAllSeries = () => toggleAllSeries(false);

// チャートのオプション
const chartOptions = reactive({
  chart: {
    type: 'line',
    stacked: true,
    toolbar: {
      show: true
    },
    zoom: {
      enabled: false,
    }
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '55%',
      endingShape: 'rounded',
      dataLabels: {
        position: 'center',
      }
    }
  },
  dataLabels: {
    enabled: true,
    enabledOnSeries: series.value.filter(s => s.type === 'bar').map((s,i) => i),
    formatter: function (val: number) {
      return val.toFixed(1) +"h"
    },
    style: {
      colors: ['#000'],
      fontSize: '10px',
    },
    dropShadow: {
      enabled: true,
      left: 2,
      top: 2,
      opacity: 0.5
    },
    background: {
      enabled: false,
    },
    textAnchor: 'middle',
  },
  stroke: {
    show: true,
    width: 1,
    curve: 'straight', // 線をスムーズに
  },
  xaxis: {
    // type: 'datetime', // datetime
    categories: xLabels.value,
    tickAmount: xLabels.value.length - 1, // または必要な数
    tickPlacement: 'on',
    labels: {
      rotate: -90, // ラベルを縦に表示（90度回転）
      rotateAlways: true,
      style: {
        fontSize: '12px'
      },
      formatter: function (val: number) {
        const date = new Date(val);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');

        if (timeFilter.value === 'day') {
          return `${year}-${month}-${day}日`;
        } else if (timeFilter.value === 'week') {
          return `${year}-${month}-${day}週`;
        } else if (timeFilter.value === 'month') {
          return `${year}年${month}月`;
        }
        return val; // デフォルト値を返す（エラー回避のため）
      }
    }
  },
  // 追加: 折れ線グラフのマーカー設定
  markers: {
    size: 4,
    colors: ['#45FF60', '#FF4560'],
    strokeWidth: 2,
    hover: {
      size: 7,
    }
  },
  // 以下は既存のchartOptionsに追加する内容
  yaxis: {
    title: {
      text: '値',
    },
    labels: {
      formatter: function(val: number) {
        return val.toFixed(1) + "h"; // 小数点以下1桁で表示し、単位「h」を追加
      }
    }
  },
  fill: {
    opacity: 1
  },
  tooltip: {
    y: {
      formatter: function (val: number) {
        return val.toFixed(1) + "h"
      }
    }
  },
  legend: {
    position: 'top',
    markers: {
      shape: series.value.map((v, index) => { return v.type === 'bar' ? 'square' : 'line' }),
    }
  },
  events: {
    legendClick: function (chartContext, seriesIndex) {
      // ここではなく、@legendClickイベントで処理するため空にしておく
    }
  }
})

// 日付をHTML input[type=date]用にフォーマットする関数
const formatDateForInput = (timestamp: number): string => {
  return dayjs(timestamp).format('YYYY-MM-DD');
};

// 開始日変更時の処理
const handleStartDateChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.value) {
    // 選択された日付をタイムスタンプに変換
    const newStartDate = dayjs(target.value).valueOf();
    // 開始日を更新
    updateStartDate(newStartDate);
    updateChart();
  }
};

// 更新時の処理
const updateChart = async () => {
  await refreshData()
  // フィルター変更時の処理（必要に応じて追加）
  if (chartRef.value?.chart) {
    // 表示するシリーズの配列を作成
    // 棒グラフのシリーズをチェックボックスに応じて振り分け
    // series.value.forEach((v, index) => {
    //   if (selectedSeries[index]) {
    //     chartRef.value.chart.showSeries(v.name);
    //   } else {
    //     chartRef.value.chart.hideSeries(v.name);
    //   }
    // });
    chartRef.value.chart.updateOptions({
      xaxis: {
        categories: xLabels.value,
        labels: {
          rotate: -90, // ラベルを縦に表示（90度回転）
          rotateAlways: true,
          style: {
            fontSize: '12px'
          },
          formatter: function (val: number) {
            const date = new Date(val);
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');

            if (timeFilter.value === 'day') {
              return `${year}-${month}-${day}日`;
            } else if (timeFilter.value === 'week') {
              return `${year}-${month}-${day}週`;
            } else if (timeFilter.value === 'month') {
              return `${year}年${month}月`;
            }
            return val; // デフォルト値を返す（エラー回避のため）
          }
        }
      }
    })

  }
}

// チャート側でレジェンドをクリックしたときのイベントハンドラ
const handleLegendClick = (chartContext: any, seriesIndex: number) => {
  selectedSeries.value[seriesIndex] = !selectedSeries.value[seriesIndex];
};

</script>


<style scoped>
/* Ensure the component takes up the full height of its parent */
:deep(.apexcharts-canvas) {
  height: 100% !important;
}

:deep(.apexcharts-graphical) {
  height: 100% !important;
}

.graph-container {
  display: flex;
  position: relative;
  width: 100%;
  height: 100vh; /* Use viewport height to ensure full browser height */
}

.toggle-filter-btn {
  position: absolute;
  left: 60px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 10;
  width: 20px;
  height: 40px;
  border: none;
  background-color: #f0f0f0;
  cursor: pointer;
  border-radius: 0 4px 4px 0;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  transition: left 0.3s ease, right 0.3s ease, border-radius 0.3s ease;
}

/* When filter panel is visible (not hidden), position the button on the right side of the filter panel */
.graph-container.filter-open .toggle-filter-btn {
  position: absolute;
  left: 250px;
  right: auto;
  border-radius: 0 4px 4px 0;
}

.filter-panel {
  flex: 0 0 250px; /* flex-grow: 0, flex-shrink: 0, flex-basis: 250px */
  width: 250px;
  padding: 15px;
  background-color: #f8f8f8;
  border-right: 1px solid #ddd;
  overflow-y: auto;
  height: 100%;
  position: relative;
}

.filter-panel.hidden {
  flex: 0 0 60px;
  width: 60px;
  overflow: visible;
  position: relative;
}

/* Styles for elements in the minimized filter panel */
.filter-panel.hidden h3 {
  writing-mode: vertical-rl;
  text-align: center;
  margin: 10px auto;
  white-space: nowrap;
}

.filter-panel.hidden .filter-item {
  position: relative;
  margin: 30px 0;
  cursor: pointer;
}

.filter-panel.hidden .filter-item > label {
  writing-mode: vertical-rl;
  text-align: center;
  margin: 10px auto;
  white-space: nowrap;
  position: relative;
}

.filter-panel.hidden .filter-item > label::after {
  content: "▶";
  position: absolute;
  top: -15px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: #666;
}

/* Styles for hover content in minimized mode */
.hover-content {
  display: none; /* Hide by default */
}

/* Show hover content when filter item is hovered in minimized mode */
.filter-panel.hidden .filter-item:hover .hover-content {
  display: block;
  position: absolute;
  left: 60px;
  top: 50%;
  transform: translateY(-50%);
  background: white;
  border: 2px solid #4a86e8;
  border-radius: 6px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
  z-index: 110;
  padding: 8px 8px 8px 20px;
  min-width: 200px;
}

/* Create a bridge between filter item and hover content to prevent gap */
.filter-panel.hidden .filter-item:hover .hover-content::before {
  content: "";
  position: absolute;
  top: 0;
  left: -60px; /* Extend to the left to bridge the entire gap */
  width: 60px; /* Width of the bridge to match the gap */
  height: 100%;
  background: transparent; /* Invisible bridge */
  z-index: 100;
}

/* Specific styling for radio groups in hover content */
.filter-panel.hidden .filter-item:hover .hover-content .radio-group {
  display: flex;
}

/* Specific styling for date input to ensure calendar icon position */
.filter-panel.hidden .hover-content input[type="date"] {
  width: 200px; /* Fixed width to ensure consistent layout */
  box-sizing: border-box; /* Include padding and border in the width */
  padding-right: 30px; /* Space for the calendar icon */
}


.filter-item {
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
}


.checkbox-item input[type="checkbox"] {
  margin-right: 8px;
}

.checkbox-item label {
  margin: 0;
}
.chart-area {
  flex-grow: 1;
  padding: 15px;
  transition: width 0.3s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.chart-area.full-width {
  width: calc(100% - 60px);
  margin-left: 60px;
}

.filter-item {
  margin-bottom: 16px;
}

.filter-item label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
}

.radio-group {
  display: flex;
  gap: 12px;
  margin-top: 4px;
  justify-content: center;
}

.custom-radio {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.custom-radio input[type="radio"] {
  position: absolute;
  opacity: 0;
}

.radio-label {
  position: relative;
  padding: 6px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  transition: all 0.2s ease;
}

.custom-radio input[type="radio"]:checked + .radio-label {
  background-color: #4a86e8;
  color: white;
  border-color: #4a86e8;
}

.custom-radio input[type="radio"]:focus + .radio-label {
  box-shadow: 0 0 0 2px rgba(74, 134, 232, 0.3);
}

.custom-radio input[type="checkbox"]:checked + .radio-label {
  background-color: #4a86e8;
  color: white;
  border-color: #4a86e8;
}

.custom-radio input[type="checkbox"]:focus + .radio-label {
  box-shadow: 0 0 0 2px rgba(74, 134, 232, 0.3);
}

.custom-radio:hover .radio-label {
  border-color: #bbb;
}

.hidden-checkbox {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}
.legend-buttons {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
  padding: 0 10px;
}

.legend-btn {
  background-color: #f0f0f0;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 5px 10px;
  margin-left: 8px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s ease;
}

.legend-btn:hover {
  background-color: #e0e0e0;
  border-color: #ccc;
}

.legend-btn:active {
  background-color: #d0d0d0;
}

.reload-btn {
  display: flex;
  align-items: center;
  margin-left: auto;
  color: inherit;
  text-decoration: none;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.reload-btn:hover {
  background-color: #f0f0f0;
}

.reload-btn .material-symbols-outlined {
  font-size: 16px;
  margin-right: 4px;
}

.reload-btn .text {
  font-size: 12px;
}

</style>
