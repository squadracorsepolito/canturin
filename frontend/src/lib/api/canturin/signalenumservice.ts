// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function AddValue(entityID: string): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1459823736, entityID) as any;
    return $resultPromise;
}

export function Create(req: $models.CreateSignalEnumReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4055206348, req) as any;
    return $resultPromise;
}

export function Delete(entityID: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1564009523, entityID) as any;
    return $resultPromise;
}

export function Get(entityID: string): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(62675348, entityID) as any;
    return $resultPromise;
}

export function GetInvalidNames(entityID: string): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3908247013, entityID) as any;
    return $resultPromise;
}

export function ListBase(): Promise<$models.BaseEntity[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(535410869) as any;
    return $resultPromise;
}

export function RemoveValues(entityID: string, req: $models.RemoveValuesReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3900736342, entityID, req) as any;
    return $resultPromise;
}

export function ReorderValue(entityID: string, req: $models.ReorderValueReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1373905824, entityID, req) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, req: $models.UpdateDescReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3205892276, entityID, req) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, req: $models.UpdateNameReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3379144738, entityID, req) as any;
    return $resultPromise;
}

export function UpdateValueDesc(entityID: string, req: $models.UpdateValueDescReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(388233975, entityID, req) as any;
    return $resultPromise;
}

export function UpdateValueIndex(entityID: string, req: $models.UpdateValueIndexReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2877683578, entityID, req) as any;
    return $resultPromise;
}

export function UpdateValueName(entityID: string, req: $models.UpdateValueNameReq): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(919586037, entityID, req) as any;
    return $resultPromise;
}
