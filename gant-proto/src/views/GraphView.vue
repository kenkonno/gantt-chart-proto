<template>
  <div class="graph-container">
    <!-- トグルボタンを追加 -->
    <button class="toggle-filter-btn" @click="toggleFilterPanel">
      {{ filterVisible ? '◀' : '▶' }}
    </button>

    <!-- フィルターパネル（左側） - v-showで表示・非表示を切り替え -->
    <div class="filter-panel" :class="{ 'hidden': !filterVisible }">
      <h3>フィルター</h3>

      <div class="filter-item">
        <label>データ選択：</label>
        <div v-for="(series, index) in series" :key="index" class="checkbox-item">
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

    <!-- グラフ表示エリア - フィルター非表示時に全幅表示 -->
    <div class="chart-area" :class="{ 'full-width': !filterVisible }">
      <apex-charts
          ref="chartRef"
          :options="chartOptions"
          :series="series"
          height="350"
          @legendClick="handleLegendClick"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, reactive, ref, Ref} from 'vue'
import ApexCharts from 'vue3-apexcharts'
import {usePileUpGraph} from "@/composable/pileUpGraph";

const {
  facilities,
  departments,
  pileUps,
  facilityTypes,
  timeFilter,
  series,
  xLabels,
  selectedSeries
} = await usePileUpGraph()

// チャートの参照を作成
const chartRef = ref<null | any>(null) as Ref<typeof ApexCharts>;

// フィルターパネルの表示状態を管理するリアクティブな変数
const filterVisible = ref(true);

// フィルターパネルの表示/非表示を切り替える関数
const toggleFilterPanel = () => {
  filterVisible.value = !filterVisible.value;
};


// 選択状態の管理
const selectedPeriod = ref('month') // 初期値は1ヶ月

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
    curve: 'straight', // 線をスムーズに
    colors: ['transparent', 'transparent', 'transparent', '#FF4560'] // 最後の値は折れ線グラフの色
  },
  xaxis: {
    type: 'category', // datetime から category に変更
    categories: computed(() => {
      return xLabels.value // 横軸のラベル
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
    colors: ['#FF4560','#45FF60'],
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
  },
  events: {
    legendClick: function (chartContext, seriesIndex) {
      // ここではなく、@legendClickイベントで処理するため空にしておく
    }
  }
})

// 更新時の処理
const updateChart = () => {
  // フィルター変更時の処理（必要に応じて追加）
  if (chartRef.value?.chart) {
    // 表示するシリーズの配列を作成
    // 棒グラフのシリーズをチェックボックスに応じて振り分け
    series.value.forEach((v, index) => {
      if (selectedSeries[index]) {
        chartRef.value.chart.showSeries(v.name);
      } else {
        chartRef.value.chart.hideSeries(v.name);
      }
    });
  }
}

// ApexCharts のイベントハンドラーのための型を定義
type ChartContext = {
  w: {
    globals: {
      collapsedSeriesIndices: number[];
      // 他に必要な型情報
    }
  };
  // その他の必要なプロパティ
};

// チャート側でレジェンドをクリックしたときのイベントハンドラ
const handleLegendClick = (chartContext: ChartContext, seriesIndex: number) => {
  console.log('レジェンドがクリックされました:', seriesIndex);

  selectedSeries[seriesIndex] = chartContext.w.globals.collapsedSeriesIndices.indexOf(seriesIndex) !== -1;
};

</script>


<style scoped>
.graph-container {
  display: flex;
  position: relative;
  width: 100%;
  height: 100%;
}

.toggle-filter-btn {
  position: absolute;
  left: 0;
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
}

.filter-panel {
  width: 250px;
  padding: 15px;
  background-color: #f8f8f8;
  border-right: 1px solid #ddd;
  transition: transform 0.3s ease;
  overflow-y: auto;
  height: 100%;
}

.filter-panel.hidden {
  transform: translateX(-100%);
  position: absolute;
}

.filter-item {
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
}

/* 期間選択用の特別なスタイル */
.period-item {
  align-items: center; /* 中央揃え */
}

.filter-label {
  margin-bottom: 8px;
  font-weight: bold;
  text-align: left;
}

/* 期間選択のラベル用 */
.period-label {
  text-align: center; /* 中央揃え */
}

.checkbox-container {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.checkbox-item {
  margin: 5px 0;
  display: flex;
  align-items: center;
  text-align: left;
  width: 100%;
}

.checkbox-item input[type="checkbox"] {
  margin-right: 8px;
}

.checkbox-item label {
  margin: 0;
}

/* 期間選択のセレクトボックス */
.period-select {
  width: auto;
  margin: 0 auto; /* 中央揃え */
}

.chart-area {
  flex-grow: 1;
  padding: 15px;
  transition: width 0.3s ease;
}

.chart-area.full-width {
  width: 100%;
}
</style>
