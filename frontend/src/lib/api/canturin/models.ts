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
    "nodes": Node[] | null;
}

export interface BusStub {
    "entityId": string;
    "name": string;
    "nodes": NodeStub[] | null;
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
    "receivers": Node[] | null;
}

export interface MessageStub {
    "entityId": string;
    "name": string;
    "signals": $internal.entityStub[] | null;
}

export interface Network {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "buses": Bus[] | null;
}

export interface NetworkStub {
    "entityId": string;
    "name": string;
    "buses": BusStub[] | null;
    "signalUnits": SignalUnitStub[] | null;
    "signalTypes": SignalTypeStub[] | null;
}

export interface Node {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "sendedMessages": Message[] | null;
}

export interface NodeStub {
    "entityId": string;
    "name": string;
    "sendedMessages": MessageStub[] | null;
}

export interface Signal {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "kind": acmelib$0.SignalKind;
    "startPos": number;
    "size": number;
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
    "kind": acmelib$0.SignalTypeKind;
    "min": number;
    "max": number;
    "scale": number;
    "offset": number;
    "referenceCount": number;
    "references": SignalReference[] | null;
}

export interface SignalTypeStub {
    "entityId": string;
    "name": string;
}

export interface SignalUnit {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "symbol": string;
    "referenceCount": number;
    "references": SignalReference[] | null;
}

export interface SignalUnitStub {
    "entityId": string;
    "name": string;
}
