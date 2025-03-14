<template>
  <div class="graph-container">
    <!-- フィルターパネル（左側） -->
    <div class="filter-panel">
      <h3>フィルター</h3>

      <div class="filter-item">
        <label>データ選択：</label>
        <div v-for="(series, index) in availableSeries" :key="index" class="checkbox-item">
          <input
              type="checkbox"
              :id="`series-${index}`"
              v-model="selectedSeries[index]"
              @change="updateChart"
          />
          <label :for="`series-${index}`">{{ series.name }}</label>
        </div>
      </div>

      <div class="filter-item">
        <label>期間選択：</label>
        <select v-model="selectedPeriod" @change="updateChart">
          <option value="week">直近1週間</option>
          <option value="month">直近1ヶ月</option>
          <option value="quarter">直近3ヶ月</option>
        </select>
      </div>
    </div>

    <!-- グラフ表示エリア -->
    <div class="chart-area">
      <apex-charts
          type="bar"
          :options="chartOptions"
          :series="filteredSeries"
          height="350"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, reactive, computed, onMounted} from 'vue'
import ApexCharts from 'vue3-apexcharts'
import {usePileUpGraph} from "@/composable/pileUpGraph";

const {
  facilities,
  departments,
  pileUps,
  facilityTypes,
  timeFilter,
  barSeries,
  lineSeries,
  xLabels
} = await usePileUpGraph()


// 日付データの生成（過去30日分）
const generateDates = (days: number) => {
  const dates = []
  const today = new Date()

  for (let i = days - 1; i >= 0; i--) {
    const date = new Date()
    date.setDate(today.getDate() - i)
    dates.push(date.toISOString().split('T')[0]) // YYYY-MM-DD形式
  }

  return dates
}

// 利用可能なシリーズデータ
const availableSeries = ref(barSeries)

// 折れ線グラフのデータを計算する関数
const calculateLineSeries = () => {
  return lineSeries
}

// 選択状態の管理
const selectedSeries = reactive([true, true, true]) // 初期値は全て表示
const selectedPeriod = ref('month') // 初期値は1ヶ月

// フィルターに基づいたデータを計算
const filteredSeries = computed(() => {
  // 期間に応じたデータの抽出
  let periodDays = 30 // デフォルト1ヶ月
  if (selectedPeriod.value === 'week') {
    periodDays = 7
  } else if (selectedPeriod.value === 'quarter') {
    periodDays = 90
  }

  // 期間に応じた日付の更新
  const filteredDates = generateDates(periodDays)

  // 表示する系列を選択
  const barSeries = availableSeries.value
      .map((series, index) => {
        if (!selectedSeries[index]) return null

        return {
          name: series.name,
          data: series.data.slice(-periodDays), // 指定期間のみ抽出
          color: series.color
        }
      })
      .filter(item => item !== null) // 非選択の項目を除外

  // 折れ線グラフのデータを計算
  const lineSeries = calculateLineSeries()

  // 棒グラフと折れ線グラフを結合
  return [...barSeries, ...lineSeries]
})

// チャートのオプション
const chartOptions = reactive({
  chart: {
    type: 'bar',
    height: 350,
    stacked: true,
    toolbar: {
      show: true
    },
    zoom: {
      enabled: true
    }
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '55%',
      endingShape: 'rounded'
    }
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    show: true,
    width: [2, 2, 2, 4], // 最後の値は折れ線グラフの線の太さ
    curve: 'smooth', // 線をスムーズに
    colors: ['transparent', 'transparent', 'transparent', '#FF4560'] // 最後の値は折れ線グラフの色
  },
  xaxis: {
    type: 'category', // datetime から category に変更
    categories: computed(() => {
      return xLabels // 横軸のラベル
    }),
    labels: {
      rotate: -90, // ラベルを縦に表示（90度回転）
      rotateAlways: true,
      style: {
        fontSize: '12px'
      },
    }
  },
  // 追加: 折れ線グラフのマーカー設定
  markers: {
    size: 4,
    colors: ['#FF4560'],
    strokeWidth: 2,
    hover: {
      size: 7,
    }
  },
  // 以下は既存のchartOptionsに追加する内容
  yaxis: {
    title: {
      text: '値'
    }
  },
  fill: {
    opacity: 1
  },
  tooltip: {
    y: {
      formatter: function (val) {
        return val + " ユニット"
      }
    }
  },
  legend: {
    position: 'top'
  }
})

// 更新時の処理
const updateChart = () => {
  // フィルター変更時の処理（必要に応じて追加）
}
</script>

<style scoped>
.graph-container {
  display: flex;
  gap: 20px;
  padding: 20px;
}

.filter-panel {
  width: 250px;
  padding: 15px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.chart-area {
  flex: 1;
  background-color: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.filter-item {
  margin-bottom: 15px;
}

.checkbox-item {
  margin: 5px 0;
  display: flex;
  align-items: center;
}

h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #333;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
}

select {
  width: 100%;
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #ddd;
}
</style>