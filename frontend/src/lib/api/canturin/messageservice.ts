// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function Get(entityID: string): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(383147390, entityID) as any;
    return $resultPromise;
}

export function GetInvalidCANIDs(entityID: string, busEntityID: string): Promise<number[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1875989265, entityID, busEntityID) as any;
    return $resultPromise;
}

export function GetInvalidMessageIDs(entityID: string): Promise<number[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2352767286, entityID) as any;
    return $resultPromise;
}

export function GetInvalidNames(entityID: string): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1897884935, entityID) as any;
    return $resultPromise;
}

export function ListBase(): Promise<$models.BaseEntity[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(417424819) as any;
    return $resultPromise;
}

export function UpdateCycleTime(entityID: string, req: $models.UpdateCycleTimeReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2224820826, entityID, req) as any;
    return $resultPromise;
}

export function UpdateDelayTime(entityID: string, req: $models.UpdateDelayTimeReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2945565099, entityID, req) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, req: $models.UpdateDescReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(866059846, entityID, req) as any;
    return $resultPromise;
}

export function UpdateMessageID(entityID: string, req: $models.UpdateMessageIDReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3298149343, entityID, req) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, req: $models.UpdateNameReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1372204560, entityID, req) as any;
    return $resultPromise;
}

export function UpdateSendType(entityID: string, req: $models.UpdateSendTypeReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1863973285, entityID, req) as any;
    return $resultPromise;
}

export function UpdateStartDelayTime(entityID: string, req: $models.UpdateStartDelayTimeReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3149167501, entityID, req) as any;
    return $resultPromise;
}

export function UpdateStaticCANID(entityID: string, req: $models.UpdateStaticCANIDReq): Promise<$models.Message> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2727991396, entityID, req) as any;
    return $resultPromise;
}
