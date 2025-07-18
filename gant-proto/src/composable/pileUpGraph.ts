import {Api} from "@/api/axios";
import {computed, reactive, ref} from "vue";
import {FacilityType} from "@/const/common";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import dayjs from "dayjs";
import {DefaultValidIndexUsers, Department, GetDefaultPileUpsResponse, User} from "@/api";

export type Series = {
    name: string;
    type: string;
    data: number[];
    color: string;
    hidden: boolean;
    highUtilizationPoints: boolean[]; // 各時点の稼働率が125%を超えているかどうかのフラグ
}

const dailyDurationOptions = [
    {key: "2週間", value: 14},
    {key: "1ヶ月", value: 30}
]
const weeklyDurationOptions = [
    {key: "1四半期", value: 13},
]
const monthlyDurationOptions = [
    {key: "年間", value: 12},
]


// ユーザー一覧。特別ref系は必要ない。
export async function usePileUpGraph() {

    // TODO: 日付の初期値がおかしい？基本データは詰まってくるはずなのであまり気にしない？

    // 設備一覧
    const {data: facilities} = reactive(await Api.getFacilities())
    // 部署一覧
    const {data: departments} = reactive(await Api.getDepartments())
    // ユーザー一覧
    const {data: userList} = reactive(await Api.getUsers())
    // フィルタ 設備の種類を定数から取得し、リアクティブに設定する
    const facilityTypes = reactive<string[]>([FacilityType.Ordered, FacilityType.Prepared])
    // 日次・週次・月次の絞り込みの選択肢
    const timeFilter = ref<DisplayType>('week');
    // 期間の選択
    const durationFilter = ref<number>(13);
    // 開始日の選択（デフォルトは本日）
    const startDate = ref<number>(normalizeDate(Date.now()));
    // 期間のオプション
    const durationOptions = computed(() => {
        switch (timeFilter.value) {
            case 'day':
                return dailyDurationOptions
            case 'week':
                return weeklyDurationOptions
            case 'month':
                return monthlyDurationOptions
            default:
                return []
        }
    })
    const updateDuration = () => {
        switch (timeFilter.value) {
            case 'day':
                return durationFilter.value = dailyDurationOptions[0].value
            case 'week':
                return durationFilter.value = weeklyDurationOptions[0].value
            case 'month':
                return durationFilter.value = monthlyDurationOptions[0].value
        }
    }

    // 積み上げ（getDefaultPileUps）
    const pileUps = reactive<GetDefaultPileUpsResponse>({
        defaultPileUps: [],
        defaultValidUserIndexes: [],
        globalStartDate: ''
    })
    let globalStartDate = ''

    const refreshData = async () => {
        const {data} = await Api.getDefaultPileUps(-1, true, facilityTypes)
        pileUps.defaultPileUps = data.defaultPileUps
        pileUps.defaultValidUserIndexes = data.defaultValidUserIndexes.map(v => {
            if (v.isHoliday) v.UserIds.length = 0
            return v
        })
        globalStartDate = data.globalStartDate
        return data
    }

    // 初回データ取得
    await refreshData()

    const labelLength = computed(() => {
        return splitByTimeFilter(pileUps.defaultValidUserIndexes.map(v => v.ValidIndex), timeFilter.value, globalStartDate).length
    })


    // 横軸のラベル作成
    const xLabels = computed(() => {
        const allXLabels = getXLabels(globalStartDate, timeFilter.value, labelLength.value)
        return filterXLabelsByDuration(allXLabels, durationFilter.value, startDate.value)
    })

    // 部署のフィルタ
    const hideSeries = ref<boolean[]>(Array.from({length: pileUps.defaultPileUps.length + 2}, () => false))

    // seriesの作成
    const series = computed(() => {
        // 部署ごとの100%基準値（各時点での稼働可能な人数 * 8時間）を計算
        const averageLineByDepartment = getAverageLine(pileUps.defaultValidUserIndexes, departments.list.map(v => v.id!), userList.list);

        // 表示する部署のIDを取得
        const displayDepartmentIds = pileUps.defaultPileUps.map(v => v.departmentId!).filter((v, i) => !hideSeries.value[i]);

        // 表示する部署の稼働可能時間の合計を計算
        const averageLabels = sumArrays(averageLineByDepartment.filter(v => displayDepartmentIds.includes(v.departmentId)).map(v => v.labels));

        // 合計稼働可能時間を時間フィルタに応じて分割
        const totalAverageData = splitByTimeFilter(averageLabels, timeFilter.value, globalStartDate);

        // 各部署のデータをパーセント表示に変換
        const barSeries = pileUps.defaultPileUps.map((v, index) => {
            // 対応する部署の100%基準値を取得
            const departmentId = v.departmentId!;
            const departmentAverage = averageLineByDepartment.find(avg => avg.departmentId === departmentId)!;

            // 部署の各時点のデータを取得
            const rawData = splitByTimeFilter(v.labels, timeFilter.value, globalStartDate);
            const departmentAverageData = splitByTimeFilter(departmentAverage.labels, timeFilter.value, globalStartDate);

            // 全部署の合計稼働可能時間を基準として、パーセント表示に変換
            const percentData = rawData.map((value, i) => {
                if (totalAverageData && i < totalAverageData.length && totalAverageData[i] > 0) {
                    // 100%基準値 = 全部署の稼働可能な人数の合計 * 8時間
                    const baseValue = totalAverageData[i] * 8;
                    // パーセント計算（小数点以下2桁まで）
                    return Number(((value / baseValue) * 100).toFixed(2));
                }
                return 0; // 基準値がない場合は0%とする
            });

            // 部署単体での稼働率を計算し、125%を超えているかどうかをチェック
            const highUtilizationPoints = rawData.map((value, i) => {
                if (departmentAverageData && i < departmentAverageData.length && departmentAverageData[i] > 0) {
                    // 100%基準値 = その部署の稼働可能な人数 * 8時間
                    const baseValue = departmentAverageData[i] * 8;
                    // パーセント計算
                    const percent = (value / baseValue) * 100;
                    // 125%を超えているかどうかをチェック
                    return percent > 125;
                }
                return false; // 基準値がない場合はfalse
            });

            return <Series>{
                color: departments.list.find(dept => dept.id === departmentId)?.color || '#CCCCCC', // 部署マスタの色を使用、見つからない場合はデフォルト色
                data: percentData,
                name: getDepartmentName(departmentId, departments.list),
                type: "bar",
                hidden: hideSeries.value[index],
                highUtilizationPoints: highUtilizationPoints // 高稼働率ポイントの情報を追加
            }
        });

        // ダミーデータ,100%線, 125%線。部署ごとにその日に稼動可能な人数を足し上げて計算する
        const lineSeries: Series[] = [{
            color: "green",
            data: Array(splitByTimeFilter(averageLabels, timeFilter.value, globalStartDate).length).fill(100),
            name: "100%",
            type: "line",
            hidden: false,
            highUtilizationPoints: [],
        }, {
            color: "red",
            data: Array(splitByTimeFilter(averageLabels, timeFilter.value, globalStartDate).length).fill(125),
            name: "125%",
            type: "line",
            hidden: false,
            highUtilizationPoints: [],
        }]

        const xLabels = getXLabels(globalStartDate, timeFilter.value, labelLength.value)

        return [...(filterSeriesByDuration(barSeries, xLabels, durationFilter.value, startDate.value)), ...(filterSeriesByDuration(lineSeries, xLabels, durationFilter.value, startDate.value))]
    })


    // 開始日を更新する関数
    const updateStartDate = (newStartDate: number) => {
        startDate.value = normalizeDate(newStartDate);
    };

    // globalStartDateのタイムスタンプ（ミリ秒）
    const globalStartTimestamp = ref(dayjs(globalStartDate).valueOf());

    return {
        facilities,
        departments,
        pileUps,
        facilityTypes,
        timeFilter,
        xLabels,
        selectedSeries: hideSeries,
        series,
        refreshData,
        durationOptions,
        durationFilter,
        updateDuration,
        startDate,
        updateStartDate,
        globalStartTimestamp
    }
}

/**
 * 複数の数値配列の要素同士を足し合わせる関数
 * @param arrays 足し合わせる数値配列の二次元配列
 * @returns 各要素を足し合わせた新しい配列
 */
const sumArrays = (arrays: number[][]): number[] => {
    // 配列が空の場合は空配列を返す
    if (arrays.length === 0) {
        return [];
    }

    // 最初の配列をベースとする
    const result = [...arrays[0]];

    // 2つ目以降の配列を足していく
    for (let i = 1; i < arrays.length; i++) {
        const currentArray = arrays[i];

        // 各要素を足していく
        for (let j = 0; j < currentArray.length; j++) {
            // resultの長さよりもインデックスが大きい場合は、resultを拡張
            if (j >= result.length) {
                result.push(currentArray[j]);
            } else {
                result[j] += currentArray[j];
            }
        }
    }

    return result;
};

/**
 * 部署ごとに稼働可能な人数を集計する関数
 * @param defaultValidUserIndexes validIndexごとに稼動可能なユーザーIDの配列
 * @param departmentIds 集計対象の部署ID配列
 * @param userList ユーザーリスト
 * @return {Array<{departmentId: number, labels: number[]}>} 部署IDと各validIndexごとの稼働可能人数の配列
 */
const getAverageLine = (defaultValidUserIndexes: DefaultValidIndexUsers[], departmentIds: number[], userList: User[]) => {
    // 結果を格納する配列
    const result: Array<{ departmentId: number, labels: number[] }> = [];

    // 各部署IDに対して処理
    departmentIds.forEach(departmentId => {
        // この部署に所属するユーザーのIDを抽出
        const departmentUserIds = userList
            .filter(user => user.department_id === departmentId)
            .map(user => user.id)
            .filter((id): id is number => id !== null && id !== undefined);

        // 各validIndexごとの稼働可能人数を格納する配列
        const labels: number[] = [];

        defaultValidUserIndexes.forEach((validIndexData) => {
            // この部署のユーザーで、現在のvalidIndexで稼働可能なユーザー数をカウント
            const count = departmentUserIds.filter(userId =>
                validIndexData.UserIds.includes(userId)
            ).length;

            labels.push(count);

        })
        // 結果に追加
        result.push({
            departmentId,
            labels
        });
    });

    return result;
};

const getXLabels = (globalStartDate: string, timeFilter: DisplayType, len: number) => {
    // DisplayTypeの特定のキーだけを選んで新しい型を作成
    type TimeFilterType = Exclude<DisplayType, 'hour'>;

    // 設定オブジェクトの型定義
    type TimeFilterConfig = {
        [key in TimeFilterType]: {
            adjustStart: (date: dayjs.Dayjs) => dayjs.Dayjs;
            addUnit: key;
        };
    };


    // 各timeFilterに対する設定を定義
    const timeFilterConfig: TimeFilterConfig = {
        day: {
            adjustStart: (date: dayjs.Dayjs) => date, // 日付の場合は調整なし
            addUnit: 'day',
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
        },
        month: {
            adjustStart: (date: dayjs.Dayjs) => date.date(1), // 月の場合、月初（1日）に調整
            addUnit: 'month',
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
        return currentDate.unix() * 1000
    });

    return labels;
};


const getDepartmentName = (departmentId: number, list: Department[]) => {
    return list.find(v => v.id === departmentId)?.name ?? ""
}


// 時間のフィルター（'day', 'week', または 'month'）に応じてラベルの値を変換または集計する関数
const splitByTimeFilter = (labels: number[], timeFilter: DisplayType, globalStartDate: string): number[] => {
    // dayjs を使用して日付を操作
    const startDate = dayjs(globalStartDate).startOf('day');

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
/**
 * 日時データの時分秒を0に統一する関数
 * @param timestamp タイムスタンプ（ミリ秒）
 * @returns 時分秒が0に統一されたタイムスタンプ（ミリ秒）
 */
const normalizeDate = (timestamp: number): number => {
    return dayjs(timestamp).startOf('day').valueOf();
};

/**
 * 指定された開始日時から指定された期間分だけxLabelsデータをフィルタリングする
 * @param xLabelsData 日付のUnixタイムスタンプ配列
 * @param duration フィルタリングする期間の長さ
 * @param startDateValue 開始日時（ミリ秒）
 * @returns フィルタリングされたxLabels
 */
const filterXLabelsByDuration = (
    xLabelsData: number[],
    duration: number,
    startDateValue: number
): number[] => {
    // 開始日時を時分秒を0に統一
    const startTime = normalizeDate(startDateValue);

    // xLabelsから開始日時に最も近いインデックスを見つける
    let currentIndex = xLabelsData.findIndex(timestamp => {
        const normalizedTimestamp = normalizeDate(timestamp);
        return normalizedTimestamp >= startTime;
    });

    // 現在日時が見つからない場合は、最初のインデックスを使用
    if (currentIndex === -1) {
        currentIndex = 0;
    }

    // フィルタリング範囲の終了インデックス
    const endIndex = Math.min(currentIndex + duration, xLabelsData.length);

    // フィルタリングされたxLabels
    return xLabelsData.slice(currentIndex, endIndex);
};

/**
 * シリーズデータとxLabelsを同じ期間でフィルタリングする
 * @param seriesData フィルタリングするシリーズデータの配列
 * @param xLabelsData 日付のUnixタイムスタンプ配列
 * @param duration フィルタリングする期間の長さ
 * @returns フィルタリングされたシリーズとxLabelsを含むオブジェクト
 */
const filterSeriesByDuration = (
    seriesData: Series[],
    xLabelsData: number[],
    duration: number,
    startDateValue: number
) => {
    // xLabelsをフィルタリング
    const filteredLabels = filterXLabelsByDuration(xLabelsData, duration, startDateValue);

    // フィルタリングされた範囲のインデックスを取得
    const startIndex = xLabelsData.indexOf(filteredLabels[0]);
    const endIndex = startIndex + filteredLabels.length;

    // 同じインデックス範囲でシリーズデータをフィルタリング
    return seriesData.map(seriesItem => {
        const filteredItem = {
            ...seriesItem,
            data: seriesItem.data.slice(startIndex, endIndex)
        };
        filteredItem.highUtilizationPoints = seriesItem.highUtilizationPoints.slice(startIndex, endIndex);
        return filteredItem;
    });
};
