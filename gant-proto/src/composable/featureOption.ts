import {getFeatureOptions} from "@/composable/auth";
import {FeatureOption} from "@/const/common";


export function available(optionName: FeatureOption) {

    const featureOptions = getFeatureOptions()

    return featureOptions?.find(v => v.name == optionName && v.enabled) != undefined
}