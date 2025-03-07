// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function Create(): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1724364363) as any;
    return $resultPromise;
}

export function Delete(entityID: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3439884680, entityID) as any;
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

export function ListBase(): Promise<$models.BaseEntity[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3689944206) as any;
    return $resultPromise;
}

export function ListBrief(): Promise<$models.SignalUnitBrief[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2557497877) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, req: $models.UpdateDescReq): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3797493755, entityID, req) as any;
    return $resultPromise;
}

export function UpdateKind(entityID: string, req: $models.UpdateSignalUnitKindReq): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3852535836, entityID, req) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, req: $models.UpdateNameReq): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1637602753, entityID, req) as any;
    return $resultPromise;
}

export function UpdateSymbol(entityID: string, req: $models.UpdateSymbolReq): Promise<$models.SignalUnit> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2276999818, entityID, req) as any;
    return $resultPromise;
}
