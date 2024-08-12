import {Api} from "@/api/axios";
import {FacilitySharedLink} from "@/api";
import {ref} from "vue";

export async function useFacilitySharedLink(facilityId: number) {

    const facilitySharedLink = ref<FacilitySharedLink>({
        id: null,
        facility_id: 0,
        uuid: "",
        created_at: undefined,
        updated_at: undefined
    })
    const {data} = await Api.getFacilitySharedLinksId(facilityId)
    if (data.facilitySharedLink != undefined) {
        facilitySharedLink.value.id = data.facilitySharedLink.id
        facilitySharedLink.value.facility_id = data.facilitySharedLink.facility_id
        facilitySharedLink.value.uuid = data.facilitySharedLink.uuid
        facilitySharedLink.value.created_at = data.facilitySharedLink.created_at
        facilitySharedLink.value.updated_at = data.facilitySharedLink.updated_at
    }

    return {facilitySharedLink}

}

