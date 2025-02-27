// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT


// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as time$0 from "../time/models.js";

export interface AddSentMessageReq {
    "interfaceNumber": number;
}

export interface AddSignalReq {
    "signalKind": SignalKind;
}

export interface AttachedNode {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "id": number;
    "interfaceNumber": number;
}

export interface BaseEntity {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
}

export interface Bus {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "type": BusType;
    "baudrate": number;
    "attachedNodes": AttachedNode[] | null;
}

export interface BusLoad {
    "percentage": number;
    "messages": BusLoadMessage[] | null;
}

export interface BusLoadMessage {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "bitsPerSec": number;
    "percentage": number;
}

export enum BusType {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    BusTypeCAN2A = "CAN_2.0A",
};

export interface CreateBusReq {
    "name": string;
    "desc": string;
    "busType": BusType;
    "baudrate": number;
}

export interface CreateNodeReq {
    "name": string;
    "desc": string;
    "nodeId": number;
    "interfaceCount": number;
}

export interface CreateSignalEnumReq {
    "name": string;
    "desc": string;
    "minSize": number;
}

export interface CreateSignalTypeReq {
    "name": string;
    "desc": string;
    "kind": SignalTypeKind;
    "size": number;
    "signed": boolean;
    "min": number;
    "max": number;
    "scale": number;
    "offset": number;
}

export interface CreateSignalUnitReq {
    "name": string;
    "desc": string;
    "kind": SignalUnitKind;
    "symbol": string;
}

export interface DeleteSignalsReq {
    "signalEntityIds": string[] | null;
}

export enum EntityKind {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    EntityKindNetwork = "network",
    EntityKindBus = "bus",
    EntityKindNode = "node",
    EntityKindMessage = "message",
    EntityKindSignal = "signal",
    EntityKindSignalType = "signal-type",
    EntityKindSignalUnit = "signal-unit",
    EntityKindSignalEnum = "signal-enum",
};

export interface EntityPath {
    "kind": EntityKind;
    "entityId": string;
    "name": string;
}

export interface EnumSignal {
    "signalEnum": SignalEnumBrief;
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
    "paths": EntityPath[] | null;
    "hasStaticCANID": boolean;
    "id": number;
    "canId": number;
    "sizeByte": number;
    "availableTrailingBytes": number;
    "maxAvailableSpace": number;
    "byteOrder": MessageByteOrder;
    "cycleTime": number;
    "sendType": MessageSendType;
    "delayTime": number;
    "startDelayTime": number;
    "signals": Signal[] | null;
    "receivers": Node0[] | null;
    "senderNode": BaseEntity;
    "parentBus": BaseEntity;
}

export enum MessageByteOrder {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    MessageByteOrderLittleEndian = "little-endian",
    MessageByteOrderBigEndian = "big-endian",
};

export enum MessageSendType {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    MessageSendTypeUnset = "unset",
    MessageSendTypeCyclic = "cyclic",
    MessageSendTypeCyclicIfActive = "cyclic_if_active",
    MessageSendTypeCyclicAndTriggered = "cyclic_and_triggered",
    MessageSendTypeCyclicIfActiveAndTriggered = "cyclic_if_active_and_triggered",
};

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
    "attachedBus": BaseEntity;
    "sentMessages": BaseEntity[] | null;
    "receivedMessages": BaseEntity[] | null;
}

export interface Reference {
    "kind": ReferenceKind;
    "entityId": string;
    "name": string;
    "children": Reference[] | null;
}

export enum ReferenceKind {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    ReferenceKindBus = "bus",
    ReferenceKindNode = "node",
    ReferenceKindMessage = "message",
    ReferenceKindSignal = "signal",
};

export interface RemoveReceivedMessagesReq {
    "interfaceNumber": number;
    "messageEntityIds": string[] | null;
}

export interface RemoveSentMessagesReq {
    "interfaceNumber": number;
    "messageEntityIds": string[] | null;
}

export interface RemoveValuesReq {
    "valueEntityIds": string[] | null;
}

export interface ReorderSignalReq {
    "signalEntityId": string;
    "from": number;
    "to": number;
}

export interface ReorderValueReq {
    "valueEntityId": string;
    "from": number;
    "to": number;
}

export interface Sidebar {
    "root": SidebarItem;
}

export interface SidebarItem {
    "kind": SidebarItemKind;
    "id": string;
    "path": string;
    "name": string;
    "children": SidebarItem[] | null;
}

export enum SidebarItemKind {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    SidebarItemKindGroup = "group",
    SidebarItemKindNetwork = "network",
    SidebarItemKindBus = "bus",
    SidebarItemKindNode = "node",
    SidebarItemKindNodeInterface = "node-interface",
    SidebarItemKindMessage = "message",
    SidebarItemKindSignal = "signal",
    SidebarItemKindSignalType = "signal-type",
    SidebarItemKindSignalUnit = "signal-unit",
    SidebarItemKindSignalEnum = "signal-enum",
};

export interface Signal {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "paths": EntityPath[] | null;
    "parentMessage": BaseEntity;
    "kind": SignalKind;
    "startPos": number;
    "size": number;
    "standard": StandardSignal;
    "enum": EnumSignal;
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
    "references": Reference[] | null;
}

export interface SignalEnumBrief {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "size": number;
}

export interface SignalEnumValue {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "index": number;
}

export enum SignalKind {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    SignalKindStandard = "standard",
    SignalKindEnum = "enum",
    SignalKindMultiplexed = "multiplexed",
};

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
    "references": Reference[] | null;
}

export interface SignalTypeBrief {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "kind": SignalTypeKind;
    "size": number;
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
    "kind": SignalUnitKind;
    "symbol": string;
    "referenceCount": number;
    "references": Reference[] | null;
}

export interface SignalUnitBrief {
    "entityId": string;
    "name": string;
    "desc": string;
    "createTime": time$0.Time;
    "kind": SignalUnitKind;
}

export enum SignalUnitKind {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    SignalUnitKindCustom = "custom",
    SignalUnitKindTemperature = "temperature",
    SignalUnitKindElectrical = "electrical",
    SignalUnitKindPower = "power",
};

export interface StandardSignal {
    "signalType": SignalTypeBrief;
    "signalUnit": BaseEntity;
}

export interface UpdateAttachedBusReq {
    "busEntityId": string;
    "interfaceNumber": number;
}

export interface UpdateBaudrateReq {
    "baudrate": number;
}

export interface UpdateBusTypeReq {
    "busType": BusType;
}

export interface UpdateByteOrderReq {
    "byteOrder": MessageByteOrder;
}

export interface UpdateCycleTimeReq {
    "cycleTime": number;
}

export interface UpdateDelayTimeReq {
    "delayTime": number;
}

export interface UpdateDescReq {
    "desc": string;
}

export interface UpdateMaxReq {
    "max": number;
}

export interface UpdateMessageIDReq {
    "messageId": number;
}

export interface UpdateMinReq {
    "min": number;
}

export interface UpdateNameReq {
    "name": string;
}

export interface UpdateNodeIDReq {
    "nodeId": number;
}

export interface UpdateOffsetReq {
    "offset": number;
}

export interface UpdateScaleReq {
    "scale": number;
}

export interface UpdateSendTypeReq {
    "sendType": MessageSendType;
}

export interface UpdateSignalEnumReq {
    "signalEnumEntityId": string;
}

export interface UpdateSignalTypeReq {
    "signalTypeEntityId": string;
}

export interface UpdateSignalUnitKindReq {
    "kind": SignalUnitKind;
}

export interface UpdateSignalUnitReq {
    "signalUnitEntityId": string;
}

export interface UpdateSizeByteReq {
    "sizeByte": number;
}

export interface UpdateStartDelayTimeReq {
    "startDelayTime": number;
}

export interface UpdateStaticCANIDReq {
    "staticCanId": number;
}

export interface UpdateSymbolReq {
    "symbol": string;
}

export interface UpdateValueDescReq {
    "valueEntityId": string;
    "desc": string;
}

export interface UpdateValueIndexReq {
    "valueEntityId": string;
    "index": number;
}

export interface UpdateValueNameReq {
    "valueEntityId": string;
    "name": string;
}
