<template>
  <div class="action-menu d-flex">
    <div class="wrapper d-flex">
      <div class="justify-middle">
        <div>メニュー</div>
      </div>
    </div>
    <AccordionHorizontal class="justify-middle">
      <template v-slot:icon>
        <span class="material-symbols-outlined">menu_open</span>
      </template>
      <template v-slot:body>
        <tippy content="（開始日 又は 日後）・工数・工程・人数 が入力されている行が対象です。">
          <input type="button" class="btn btn-sm btn-outline-dark" value="リスケ（工数h）重視"
                 @click="emits('setScheduleByPersonDay')">
        </tippy>
        <tippy content="開始日・終了日・工数・工程が入力されている行が対象です。担当者が入力されている行は無視されます。">
          <input type="button" class="btn btn-sm btn-outline-dark" value="リスケ(日付)重視"
                 @click="emits('setScheduleByFromTo')">
        </tippy>
      </template>
    </AccordionHorizontal>
    <AccordionHorizontal class="justify-middle">
      <template v-slot:icon>
        <span class="material-symbols-outlined">filter_list</span>
      </template>
      <template v-slot:body>
        <div class="justify-middle">
          <div class="filter">
            <label v-for="item in ganttFacilityHeader" :key="item" class="side-menu-cell">
              <input type="checkbox" v-model="item.visible"/>{{ item.name }}
            </label>
          </div>
        </div>
      </template>
    </AccordionHorizontal>
    <div class="d-flex justify-middle">
      <div class="form-check">
        <input class="form-check-input" type="radio" name="displayType" id="byDay" value="day" v-model="displayType" @change="emits('updateDisplayType', $event.target.value)">
        <label class="form-check-label" for="byDay">
          日毎
        </label>
      </div>
      <div class="form-check">
        <input class="form-check-input" type="radio" name="displayType" id="byWeek" value="week" v-model="displayType" @change="emits('updateDisplayType', $event.target.value)">
        <label class="form-check-label" for="byWeek">
          週次
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AccordionHorizontal from "@/components/accordionHorizontal/AccordionHorizontal.vue";
import {Tippy} from "vue-tippy";
import {DisplayType, GanttFacilityHeader} from "@/composable/ganttFacilityMenu";
import {ref} from "vue";

type GanttFacilityMenuProps = {
  ganttFacilityHeader: GanttFacilityHeader[],
  displayType: DisplayType,
}

const emits = defineEmits(["setScheduleByPersonDay", "setScheduleByFromTo", "updateDisplayType"])
const props = defineProps<GanttFacilityMenuProps>()

const ganttFacilityHeader = ref(props.ganttFacilityHeader)
const displayType = ref(props.displayType)
</script>