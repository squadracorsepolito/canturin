// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function AddValue(enumEntID: string): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1459823736, enumEntID) as any;
    return $resultPromise;
}

export function Create(name: string, desc: string, minSize: number): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4055206348, name, desc, minSize) as any;
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

export function GetNames(): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3250263636) as any;
    return $resultPromise;
}

export function RemoveValues(enumEntID: string, ...valueEntIDs: string[]): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3900736342, enumEntID, valueEntIDs) as any;
    return $resultPromise;
}

export function ReorderValue(enumEntID: string, valueEntID: string, $from: number, to: number): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1373905824, enumEntID, valueEntID, $from, to) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, desc: string): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3205892276, entityID, desc) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, name: string): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3379144738, entityID, name) as any;
    return $resultPromise;
}

export function UpdateValueDesc(enumEntID: string, valueEntID: string, desc: string): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(388233975, enumEntID, valueEntID, desc) as any;
    return $resultPromise;
}

export function UpdateValueIndex(enumEntID: string, valueEntID: string, index: number): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2877683578, enumEntID, valueEntID, index) as any;
    return $resultPromise;
}

export function UpdateValueName(enumEntID: string, valueEntID: string, name: string): Promise<$models.SignalEnum> & { cancel(): void } {
    let $resultPromise = $Call.ByID(919586037, enumEntID, valueEntID, name) as any;
    return $resultPromise;
}
