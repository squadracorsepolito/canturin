// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function Create(req: $models.CreateBusReq): Promise<$models.Bus> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3897245259, req) as any;
    return $resultPromise;
}

export function Get(entityID: string): Promise<$models.Bus> & { cancel(): void } {
    let $resultPromise = $Call.ByID(920363765, entityID) as any;
    return $resultPromise;
}

export function GetInvalidNames(entityID: string): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4181252272, entityID) as any;
    return $resultPromise;
}

export function GetLoad(entityID: string): Promise<$models.BusLoad> & { cancel(): void } {
    let $resultPromise = $Call.ByID(541458873, entityID) as any;
    return $resultPromise;
}

export function ListBase(): Promise<$models.BaseEntity[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3115506830) as any;
    return $resultPromise;
}

export function UpdateBaudrate(entityID: string, req: $models.UpdateBaudrateReq): Promise<$models.Bus> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4227651016, entityID, req) as any;
    return $resultPromise;
}

export function UpdateBusType(entityID: string, req: $models.UpdateBusTypeReq): Promise<$models.Bus> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1520949396, entityID, req) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, req: $models.UpdateDescReq): Promise<$models.Bus> & { cancel(): void } {
    let $resultPromise = $Call.ByID(982331387, entityID, req) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, req: $models.UpdateNameReq): Promise<$models.Bus> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3117407681, entityID, req) as any;
    return $resultPromise;
}
