<template>
  <tr>

    <td class="side-menu-cell"></td><!-- css hack min-height -->
    <gantt-td :visible="ganttFacilityHeader[0].visible" class="text-start">
      <unit-toggle-button :is-open="false" @toggle="$emit('toggle-unit', unitId)"/>
        <span class="align-middle">{{ unitName }}</span>
    </gantt-td>
    <gantt-td :visible="ganttFacilityHeader[1].visible">‐</gantt-td>
    <gantt-td :visible="ganttFacilityHeader[2].visible">‐</gantt-td>
    <gantt-td :visible="ganttFacilityHeader[3].visible" style="min-width: 8rem;">
      <UserMultiselect
          :userList="userList"
          :ticketUser="ticketUserList"
          :disabled="true"
      />
    </gantt-td>
    <gantt-td :visible="ganttFacilityHeader[4].visible">
      <FormNumber
          class="small-numeric"
          :value="numberOfWorker"
          :min="1"
          :disabled="true"
      />
    </gantt-td>
    <gantt-td :visible="ganttFacilityHeader[5].visible">
      <FormNumber
          class="small-numeric"
          :value="estimate"
          :disabled="true"
      />
    </gantt-td>
    <gantt-td :visible="ganttFacilityHeader[6].visible">‐</gantt-td>
    <gantt-td :visible="ganttFacilityHeader[7].visible">
      <input
          type="date"
          :value="minStartDate"
          :disabled="true"
      />
    </gantt-td>
    <gantt-td :visible="ganttFacilityHeader[8].visible">
      <input
          type="date"
          :value="maxEndDate"
          :disabled="true"
      />
    </gantt-td>
    <gantt-td :visible="ganttFacilityHeader[9].visible">
      <FormNumber
          class="middle-numeric"
          :value="overallProgress"
          :disabled="true"
          :min="0"
      />
    </gantt-td>
    <gantt-td :visible="ganttFacilityHeader[10].visible" v-if="allowed('UPDATE_TICKET')">‐</gantt-td>
  </tr>
</template>

<style>
@import '@/assets/gantt-override.scss';
</style>
<style lang="scss" scoped>
@import '@/assets/gantt';
</style>

<script setup lang="ts">
import GanttTd from "@/components/gantt/GanttTd.vue";
import UserMultiselect from "@/components/form/UserMultiselect.vue";
import FormNumber from "@/components/form/FormNumber.vue";
import { GanttFacilityHeader } from "@/composable/ganttFacilityMenu";
import UnitToggleButton from "@/components/ganttFacility/UnitToggleButton.vue";
import {GanttChartGroup, GanttRow} from "@/composable/ganttFacility";
import {TicketUser} from "@/api";
import {allowed} from "@/composable/role";

// プロパティの定義
const props = defineProps({
  // ガントチャートのヘッダー設定
  ganttFacilityHeader: {
    type: Array as () => GanttFacilityHeader[],
    required: true
  },
  // チケット情報
  ganttChartGroup: {
    type: Object as () => GanttChartGroup,
    required: true
  },
  // ユーザーリスト
  userList: {
    type: Array,
    required: true
  },
  unitName: {
    type: String,
    require: true,
  },
  unitId: {
    type: Number,
    require: true,
  }
});

// イベント定義
const emit = defineEmits(["toggle-unit"]);

const ticketUserList = props.ganttChartGroup?.rows.reduce((prev: TicketUser[], current: GanttRow) => {
  return [...prev , ...(current.ticketUsers || [])]
}, [] as TicketUser[])

const numberOfWorker = props.ganttChartGroup?.rows.reduce((prev: number, current: GanttRow) => {
  return prev + (current.ticket?.number_of_worker ?? 0)
}, 0)

const estimate = props.ganttChartGroup?.rows.reduce((prev: number, current: GanttRow) => {
  return prev + (current.ticket?.estimate ?? 0)
}, 0)

// 最小開始日と最大終了日を取得
const minStartDate = props.ganttChartGroup?.rows.reduce((minDate: string | null, current: GanttRow) => {
  const currentStartDate = current.ticket?.start_date;
  if (!currentStartDate) return minDate;
  return (!minDate || currentStartDate < minDate) ? currentStartDate : minDate;
}, null) ?? null;

const maxEndDate = props.ganttChartGroup?.rows.reduce((maxDate: string | null, current: GanttRow) => {
  const currentEndDate = current.ticket?.end_date;
  if (!currentEndDate) return maxDate;
  return (!maxDate || currentEndDate > maxDate) ? currentEndDate : maxDate;
}, null) ?? null;


// 進捗率の計算
const calculateProgress = () => {
  // rowsが存在しない場合は0を返す
  if (!props.ganttChartGroup?.rows?.length) return 0;

  // 合計工数と消化工数を計算
  const { totalManHours, consumedManHours } = props.ganttChartGroup.rows.reduce((acc, row) => {
    // 工数があるかチェック（man_hoursプロパティが存在すると仮定）
    const manHours = row.ticket?.estimate ?? 0;
    // 進捗率があるかチェック（progressプロパティが存在すると仮定）
    const progress = row.ticket?.progress_percent ?? 0;

    // この行で消化した工数を計算（工数 × 進捗率）
    const rowConsumedManHours = manHours * (progress / 100);

    // 累計に加算
    return {
      totalManHours: acc.totalManHours + manHours,
      consumedManHours: acc.consumedManHours + rowConsumedManHours
    };
  }, { totalManHours: 0, consumedManHours: 0 });

  // 全体の工数が0の場合は進捗率0を返す
  if (totalManHours === 0) return 0;

  // 全体の進捗率を計算（消化工数 ÷ 全体の工数 × 100）
  const overallProgress = (consumedManHours / totalManHours) * 100;

  // 小数点以下2桁に丸める
  return Math.round(overallProgress * 100) / 100;
};

// 計算された進捗率
const overallProgress = calculateProgress();

</script>