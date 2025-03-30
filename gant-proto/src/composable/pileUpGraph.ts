import {Api} from "@/api/axios";
import {computed, reactive, ref} from "vue";
import {FacilityType} from "@/const/common";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import dayjs from "dayjs";
import {DefautValidIndexUsers, Department, User} from "@/api";

export type Series = {
    name: string;
    type: string;
    data: number[];
    color: string;
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

    // TODO: 積み上げのパーセント表示
    // TODO: 100%を超えたら強調表示する
    // TODO: 値に表示するものが謎？今の単位は時間になっているから人日にすればよさそう。

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

    // 積み上げ（getDefaultPileUps
    const {data: pileUps} = reactive(await Api.getDefaultPileUps(-1, true, facilityTypes))
    let globalStartDate: string
    const refreshData = async () => {
        const {data} = await Api.getDefaultPileUps(-1, true, facilityTypes)
        pileUps.defaultPileUps = data.defaultPileUps
        pileUps.defaultValidUserIndexes = data.defaultValidUserIndexes
        globalStartDate = data.globalStartDate
        return data
    }
    await refreshData()

    const labelLength = computed(() => {
        return splitByTimeFilter(pileUps.defaultValidUserIndexes.map(v => v.ValidIndex), timeFilter.value, globalStartDate).length
    })

    const color = ["#008FFB", "#FF4560", "#FFEA00"]

    // 横軸のラベル作成
    const xLabels = computed(() => {
        const allXLabels = getXLabels(globalStartDate, timeFilter.value, labelLength.value)
        return filterXLabelsByDuration(allXLabels, durationFilter.value)
    })

    // 部署のフィルタ
    const selectedSeries = reactive<boolean[]>(Array.from({length: pileUps.defaultPileUps.length + 2}, () => true))

    // seriesの作成
    const series = computed(() => {
        const barSeries = pileUps.defaultPileUps.map(v => {
            return <Series>{
                color: color.pop(), // TODO: 色を部署に持たせる
                data: splitByTimeFilter(v.labels, timeFilter.value, globalStartDate),
                name: getDepartmentName(v.departmentId!, departments.list),
                type: "bar"
            }
        })
        // 部署ごとの稼動可能人数 TODO: 休みの考慮がわけわからん。たぶんいらないと思うけど相談する。
        // TODO: getAverageLineで各部署のそのxAxisでの上限を取得すれば100%超過の計算が可能
        const averageLineByDepartment = getAverageLine(pileUps.defaultValidUserIndexes, departments.list.map(v => v.id!), userList.list);
        const displayDepartmentIds = pileUps.defaultPileUps.map(v => v.departmentId!).filter((v, i) => selectedSeries[i])
        const averageLabels = sumArrays(averageLineByDepartment.filter(v => displayDepartmentIds.includes(v.departmentId)).map(v => v.labels))
        // TODO: 週次表示の時に最終週が日曜日までないとちょっと違和感のある表示になる。
        // TODO: ふと思ったけど山積みは作業者だけで考えたほうが良い？
        // ダミーデータ,100%線, 125%線。部署ごとにその日に稼動可能な人数を足し上げて計算する
        const lineSeries: Series[] = [{
            color: "green",
            data: splitByTimeFilter(averageLabels.map(v => v * 8), timeFilter.value, globalStartDate),
            name: "100%",
            type: "line"
        }, {
            color: "red",
            data: splitByTimeFilter(averageLabels.map(v => v * 8 * 1.25), timeFilter.value, globalStartDate),
            name: "125%",
            type: "line"
        }]

        const xLabels = getXLabels(globalStartDate, timeFilter.value, labelLength.value)

        return [...(filterSeriesByDuration(barSeries, xLabels, durationFilter.value)), ...(filterSeriesByDuration(lineSeries, xLabels, durationFilter.value))]
    })


    return {
        facilities,
        departments,
        pileUps,
        facilityTypes,
        timeFilter,
        xLabels,
        selectedSeries,
        series,
        refreshData,
        durationOptions,
        durationFilter,
        updateDuration
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
const getAverageLine = (defaultValidUserIndexes: DefautValidIndexUsers[], departmentIds: number[], userList: User[]) => {
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

        defaultValidUserIndexes.forEach((validIndexData, i) => {
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
 * 現在日時から指定された期間分だけxLabelsデータをフィルタリングする
 * @param xLabelsData 日付のUnixタイムスタンプ配列
 * @param duration フィルタリングする期間の長さ
 * @returns フィルタリングされたxLabels
 */
const filterXLabelsByDuration = (
    xLabelsData: number[],
    duration: number
): number[] => {
    // 現在の日時を取得し、時分秒を0に統一
    const now = normalizeDate(Date.now());

    // xLabelsから現在日時に最も近いインデックスを見つける
    let currentIndex = xLabelsData.findIndex(timestamp => {
        const normalizedTimestamp = normalizeDate(timestamp);
        return normalizedTimestamp >= now;
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
    duration: number
) => {
    // xLabelsをフィルタリング
    const filteredLabels = filterXLabelsByDuration(xLabelsData, duration);

    // フィルタリングされた範囲のインデックスを取得
    const startIndex = xLabelsData.indexOf(filteredLabels[0]);
    const endIndex = startIndex + filteredLabels.length;

    // 同じインデックス範囲でシリーズデータをフィルタリング

    return seriesData.map(seriesItem => ({
        ...seriesItem,
        data: seriesItem.data.slice(startIndex, endIndex)
    }))
};
