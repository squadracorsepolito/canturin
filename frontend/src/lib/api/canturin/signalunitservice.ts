// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function Close(entityID: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3488152321, entityID) as any;
    return $resultPromise;
}

export function Get(entityID: string): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3658353397, entityID) as any;
    return $resultPromise;
}

export function GetInvalidNames(entityID: string): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(593412272, entityID) as any;
    return $resultPromise;
}

export function GetOpen(entityID: string): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2383998009, entityID) as any;
    return $resultPromise;
}

export function Open(entityID: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3212172551, entityID) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, newDesc: string): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3797493755, entityID, newDesc) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, newName: string): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1637602753, entityID, newName) as any;
    return $resultPromise;
}
