// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function Create(kind: $models.SignalTypeKind, name: string, desc: string, size: number, signed: boolean, min: number, max: number, scale: number, offset: number): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2414567789, kind, name, desc, size, signed, min, max, scale, offset) as any;
    return $resultPromise;
}

export function Get(entityID: string): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2443930791, entityID) as any;
    return $resultPromise;
}

export function GetInvalidNames(entityID: string): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1761936154, entityID) as any;
    return $resultPromise;
}

export function GetNames(): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4121102037) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, desc: string): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3391491109, entityID, desc) as any;
    return $resultPromise;
}

export function UpdateMax(entityID: string, max: number): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1547565842, entityID, max) as any;
    return $resultPromise;
}

export function UpdateMin(entityID: string, min: number): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1177178296, entityID, min) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, name: string): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(470311103, entityID, name) as any;
    return $resultPromise;
}

export function UpdateOffset(entityID: string, offset: number): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(226325631, entityID, offset) as any;
    return $resultPromise;
}

export function UpdateScale(entityID: string, scale: number): Promise<$models.SignalType> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3856836030, entityID, scale) as any;
    return $resultPromise;
}
