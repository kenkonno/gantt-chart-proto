import {Api} from "@/api/axios";
import {reactive, ref} from "vue";
import {Emit, FacilityStatus, FacilityType} from "@/const/common";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import dayjs from "dayjs";
import {Department} from "@/api";

export type Series = {
    name: string;
    type: string;
    data: number[];
    color: string;
}


// ユーザー一覧。特別ref系は必要ない。
export async function usePileUpGraph() {
    // 設備一覧
    const {data: facilities} = reactive(await Api.getFacilities())
    // 部署一覧
    const {data: departments} = reactive(await Api.getDepartments())
    // フィルタ 設備の種類を定数から取得し、リアクティブに設定する
    const facilityTypes = reactive<string[]>([FacilityType.Ordered, FacilityType.Prepared])
    // 日次・週次・月次の絞り込みの選択肢
    const timeFilter = ref<DisplayType>('day');

    // 積み上げ（getDefaultPileUps
    const {data: pileUps} = await Api.getDefaultPileUps(-1, true, facilityTypes)

    const globalStartDate = pileUps.globalStartDate

    const color = ["#008FFB", "#FF4560", "#FFEA00"]
    const barSeries = pileUps.defaultPileUps.map(v => {
        return <Series>{
            color: color.pop(),
            data: splitByTimeFilter(v.labels, timeFilter.value, globalStartDate),
            name: getDepartmentName(v.departmentId!, departments.list),
            type: "bar"
        }
    })

    // ダミーデータ,100%線, 125%線。部署ごとにその日に稼動可能な人数を足し上げて計算する
    const lineSeries: Series[] = [{
        color: "black",
        data: splitByTimeFilter(pileUps.defaultPileUps[0].labels, timeFilter.value, globalStartDate),
        name: "100%",
        type: "line"
    }, {
        color: "red",
        data: splitByTimeFilter(pileUps.defaultPileUps[1].labels, timeFilter.value, globalStartDate),
        name: "100%",
        type: "line"
    }]

    // 横軸のラベル作成
    const xLabels = getXLabels(globalStartDate, timeFilter.value, barSeries[0].data.length)

    return {facilities, departments, pileUps, facilityTypes, timeFilter, barSeries, lineSeries, xLabels}
}

const getXLabels = (globalStartDate: string, timeFilter: DisplayType, len: number) => {
    // DisplayTypeの特定のキーだけを選んで新しい型を作成
    type TimeFilterType = Exclude<DisplayType, 'hour'>;

    // 設定オブジェクトの型定義
    type TimeFilterConfig = {
        [key in TimeFilterType]: {
            adjustStart: (date: dayjs.Dayjs) => dayjs.Dayjs;
            addUnit: key;
            format: string;
        };
    };


    // 各timeFilterに対する設定を定義
    const timeFilterConfig: TimeFilterConfig = {
        day: {
            adjustStart: (date: dayjs.Dayjs) => date, // 日付の場合は調整なし
            addUnit: 'day',
            format: 'YYYY-MM-DD'
        },
        week: {
            adjustStart: (date: dayjs.Dayjs) => {
                // 週の場合、月曜日（1）に調整
                const dayOfWeek = date.day();
                // dayjsでは日曜日が0、月曜日が1...
                const daysToSubtract = dayOfWeek === 0 ? 6 : dayOfWeek - 1; // 日曜なら6日引く、それ以外はday-1
                return date.subtract(daysToSubtract, 'day');
            },
            addUnit: 'week',
            format: 'YYYY-MM-DD週'
        },
        month: {
            adjustStart: (date: dayjs.Dayjs) => date.date(1), // 月の場合、月初（1日）に調整
            addUnit: 'month',
            format: 'YYYY-MM'
        }
    };

    // 現在のフィルタの設定を取得
    const config = timeFilterConfig[timeFilter as TimeFilterType];

    // 開始日を調整
    const startDate = config.adjustStart(dayjs(globalStartDate));

    // 結果のラベル配列を生成
    const labels = Array.from({length: len}, (_, i) => {
        // i=0の場合は追加なし、それ以外はi単位追加
        const currentDate = i === 0 ? startDate : startDate.add(i, config.addUnit);

        // フォーマットによって変換
        return config.format.includes('週')
            ? `${currentDate.format('YYYY-MM-DD')}週`
            : currentDate.format(config.format);
    });

    return labels;
};


const getDepartmentName = (departmentId: number, list: Department[]) => {
    return list.find(v => v.id === departmentId)?.name ?? ""
}


// 時間のフィルター（'day', 'week', または 'month'）に応じてラベルの値を変換または集計する関数
const splitByTimeFilter = (labels: number[], timeFilter: DisplayType, globalStartDate: string): number[] => {
    // dayjs を使用して日付を操作
    const startDate = dayjs(globalStartDate);

    // 日次の場合は集計なし（各値をそのままnumber型に変換して返す）
    if (timeFilter === 'day') {
        return labels
    }

    // 週次の場合（月曜日～日曜日で集計）
    if (timeFilter === 'week') {
        const weeklyData: number[] = [];
        let weekSum = 0;

        // 最初の日の曜日を確認（0=日曜日, 1=月曜日, ..., 6=土曜日）
        const firstDayOfWeek = startDate.day();
        // 月曜日から始めるため、最初の週の開始インデックスを調整
        const adjustedWeekStart = firstDayOfWeek === 0 ? 6 : firstDayOfWeek - 1;

        // 各値を週ごとに集計
        labels.forEach((label, index) => {
            const currentDay = (index + adjustedWeekStart) % 7;
            weekSum += label;

            // 日曜日または最終要素なら週の合計を追加
            if (currentDay === 6 || index === labels.length - 1) {
                weeklyData.push(weekSum);
                weekSum = 0;
            }
        });

        return weeklyData;
    }

    // 月次の場合（1日から月末までで集計）
    if (timeFilter === 'month') {
        const monthlyData: number[] = [];
        let monthSum = 0;
        let currentMonth = startDate.month();

        // 各値を月ごとに集計
        labels.forEach((label, index) => {
            const currentDate = startDate.add(index, 'day');
            const dayMonth = currentDate.month();

            // 月が変わったら集計をリセット
            if (dayMonth !== currentMonth && index > 0) {
                monthlyData.push(monthSum);
                monthSum = 0;
                currentMonth = dayMonth;
            }

            monthSum += label;

            // 最終要素の場合も月の合計を追加
            if (index === labels.length - 1) {
                monthlyData.push(monthSum);
            }
        });

        return monthlyData;
    }

    // 未知のtimeFilterの場合は空の配列を返す
    return [];
};