// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT


// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as acmelib$0 from "../github.com/squadracorsepolito/acmelib/models.js";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as time$0 from "../time/models.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $internal from "./internal.js";

export interface Bus {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "nodeInterfaces": NodeInterface[] | null;
}

export interface BusBase {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
}

export interface History {
    "operationCount": number;
    "currentIndex": number;
}

export interface Message {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "id": acmelib$0.MessageID;
    "hasStaticCANID": boolean;
    "canId": acmelib$0.CANID;
    "sizeByte": number;
    "byteOrder": acmelib$0.MessageByteOrder;
    "signals": Signal[] | null;
    "receivers": Node0[] | null;
}

export interface Node {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "id": number;
    "interfaces": NodeInterface[] | null;
}

export interface Node0 {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "sendedMessages": Message[] | null;
}

export interface NodeInterface {
    "number": number;
    "attachedBus": BusBase;
    "sentMessages": NodeMessage[] | null;
    "receivedMessages": NodeMessage[] | null;
}

export interface NodeMessage {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
}

export interface SidebarNode {
    "kind": SidebarNodeKind;
    "name": string;
    "entityId": string;
    "parentEntityIds": string[] | null;
    "children": SidebarNode[] | null;
}

export enum SidebarNodeKind {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    SidebarNodeKindNetwork = "network",
    SidebarNodeKindBus = "bus",
    SidebarNodeKindNode = "node",
    SidebarNodeKindMessage = "message",
    SidebarNodeKindSignalType = "signal-type",
    SidebarNodeKindSignalUnit = "signal-unit",
    SidebarNodeKindSignalEnum = "signal-enum",
};

export interface Signal {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "kind": acmelib$0.SignalKind;
    "startPos": number;
    "size": number;
}

export interface SignalEnum {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "size": number;
    "minSize": number;
    "maxIndex": number;
    "values": SignalEnumValue[] | null;
    "references": SignalReference[] | null;
}

export interface SignalEnumValue {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "index": number;
}

export interface SignalReference {
    "bus": $internal.entityStub;
    "node": $internal.entityStub;
    "message": $internal.entityStub;
    "signal": $internal.entityStub;
}

export interface SignalType {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "kind": SignalTypeKind;
    "size": number;
    "signed": boolean;
    "min": number;
    "max": number;
    "scale": number;
    "offset": number;
    "referenceCount": number;
    "references": SignalReference[] | null;
}

export enum SignalTypeKind {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    SignalTypeKindCustom = "custom",
    SignalTypeKindFlag = "flag",
    SignalTypeKindInteger = "integer",
    SignalTypeKindDecimal = "decimal",
};

export interface SignalUnit {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "symbol": string;
    "referenceCount": number;
    "references": SignalReference[] | null;
}
