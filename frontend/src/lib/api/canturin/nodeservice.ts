// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function AddSentMessage(entityID: string, req: $models.AddSentMessageReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4144588181, entityID, req) as any;
    return $resultPromise;
}

export function Create(req: $models.CreateNodeReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3044972449, req) as any;
    return $resultPromise;
}

export function Delete(entityID: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3040942918, entityID) as any;
    return $resultPromise;
}

export function Get(entityID: string): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(963582531, entityID) as any;
    return $resultPromise;
}

export function GetInvalidNames(entityID: string): Promise<string[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3477642942, entityID) as any;
    return $resultPromise;
}

export function GetInvalidNodeIDs(entityID: string): Promise<number[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2336447840, entityID) as any;
    return $resultPromise;
}

export function ListBase(): Promise<$models.BaseEntity[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(720023600) as any;
    return $resultPromise;
}

export function RemoveReceivedMessages(entityID: string, req: $models.RemoveReceivedMessagesReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(507733122, entityID, req) as any;
    return $resultPromise;
}

export function RemoveSentMessages(entityID: string, req: $models.RemoveSentMessagesReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(377191259, entityID, req) as any;
    return $resultPromise;
}

export function UpdateAttachedBus(entityID: string, req: $models.UpdateAttachedBusReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2257752090, entityID, req) as any;
    return $resultPromise;
}

export function UpdateDesc(entityID: string, req: $models.UpdateDescReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(663236641, entityID, req) as any;
    return $resultPromise;
}

export function UpdateName(entityID: string, req: $models.UpdateNameReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2144136979, entityID, req) as any;
    return $resultPromise;
}

export function UpdateNodeID(entityID: string, req: $models.UpdateNodeIDReq): Promise<$models.Node> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1325739537, entityID, req) as any;
    return $resultPromise;
}
