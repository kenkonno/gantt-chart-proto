<template>
  <div class="action-menu d-flex">
    <div class="wrapper d-flex">
      <div class="justify-middle">
        <div>メニュー</div>
      </div>
    </div>
    <AccordionHorizontal class="justify-middle">
      <template v-slot:icon>
        <span class="material-symbols-outlined">filter_list</span>
      </template>
      <template v-slot:body>
        <div class="justify-middle">
          <div class="filter">
            <label v-for="item in ganttAllHeader" :key="item" class="side-menu-cell">
              <input type="checkbox" v-model="item.visible"/>{{ item.name }}
            </label>
          </div>
        </div>
      </template>
    </AccordionHorizontal>
    <div class="d-flex justify-middle">
      <div class="form-check">
        <input class="form-check-input" type="radio" name="displayType" id="byDay" v-model="displayType" value="day"
               @change="emits('updateDisplayType', $event.target.value)">
        <label class="form-check-label" for="byDay">
          日毎
        </label>
      </div>
      <div class="form-check">
        <input class="form-check-input" type="radio" name="displayType" id="byWeek" v-model="displayType"
               value="week" @change="emits('updateDisplayType', $event.target.value)">
        <label class="form-check-label" for="byWeek">
          週次
        </label>
      </div>
    </div>
  </div>
</template>
<style lang="scss" scoped>
nav {
  padding: 10px;

  > div {
    width: 100%;
    text-align: left;

    > select {
      margin: 0 5px;
    }
  }

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
<script setup lang="ts">
import AccordionHorizontal from "@/components/accordionHorizontal/AccordionHorizontal.vue";
import {ref} from "vue";
import {DisplayType, Header} from "@/composable/ganttAllMenu";

type GanttFacilityMenuProps = {
  ganttAllHeader: Header[],
  displayType: DisplayType,
}

const emits = defineEmits(["updateDisplayType"])
const props = defineProps<GanttFacilityMenuProps>()
const ganttAllHeader = ref(props.ganttAllHeader)
const displayType = ref(props.displayType)

</script>