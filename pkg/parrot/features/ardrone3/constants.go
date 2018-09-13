package ardrone3

const (
	ip = "192.168.42.1"

	// Ports
	c2dPort        = 54321 // Client to device
	d2cPort        = 43210 // Device to client
	rtpStreamPort  = 55004
	rtpControlPort = 55005
	discoveryPort  = 44444

	// libARNetworkAL/Includes/libARNetworkAL/ARNETWORKAL_Manager.h
	ARNETWORKAL_MANAGER_DEFAULT_ID_MAX uint16 = 256

	// https://developer.parrot.com/docs/bebop/ARSDK_Protocols.pdf

	// Buffers
	// (ARNETWORKAL_Frame_t identifiers)
	bufferNonack byte = 10
	// bufferAck       byte = 11
	// bufferEmergency byte = 12
	// bufferReports   byte = 127

	// Frame types
	frameTypeUninitialized  byte = 0
	frameTypeAck            byte = 1
	frameTypeData           byte = 2
	frameTypeDataLowLatency byte = 3
	frameTypeDataWithAck    byte = 4

	// https://github.com/Parrot-Developers/arsdk-xml/blob/master/xml/ardrone3.xml

	// Projects
	projectCommon   byte = 0
	projectARDrone3 byte = 1

	// Common classes
	classCommon byte = 4
	// classCommonState byte = 5

	// ARDrone3 classes
	classPiloting byte = 0
	// classPilotingState       byte = 4
	classMediaRecord byte = 7
	// classMediaRecordState    byte = 8
	// classMediaRecordEvent    byte = 3
	classSpeedSettings byte = 11
	// classSpeedSettingsState  byte = 12
	classMediaStreaming byte = 21
	// classMediaStreamingState byte = 22

	// Common commands
	cmdAllStates byte = 0

	// Common state commands
	cmdAllStatesChanged byte = 0

	// Piloting commands
	cmdFlatTrim byte = 0
	cmdTakeOff  byte = 1
	cmdPCMD     byte = 2
	cmdLanding  byte = 3
	// cmdEmergency       byte = 4
	// cmdNavigateHome    byte = 5
	// cmdAutoTakeOffMode byte = 6

	// // Piloting state commands
	// cmdFlatTrimChanged          byte = 0
	// cmdFlyingStateChanged       byte = 1
	// cmdAlertStateChanged        byte = 2
	// cmdNavigateHomeStateChanged byte = 3
	// cmdPositionChanged          byte = 4
	// cmdSpeedChanged             byte = 5
	// cmdAttitudeChanged          byte = 6
	// cmdAutoTakeOffModeChanged   byte = 7
	// cmdAltitudeChanged          byte = 8

	// // Flying states
	// flyingStateLanded           byte = 0
	// flyingStateTakingOff        byte = 1
	// flyingStateHovering         byte = 2
	// flyingStateFlying           byte = 3
	// flyingStateLanding          byte = 4
	// flyingStateEmergency        byte = 5
	// flyingStateEmergencyLanding byte = 9

	// MediaRecord commands
	cmdVideo byte = 1

	// MediaRecord commands args
	argVideoRecordStop  byte = 0
	argVideoRecordStart byte = 1

	// SpeedSettings commands
	cmdHullProtection byte = 2
	cmdOutdoor        byte = 3

	// MediaStreaming commands
	cmdVideoEnable     byte = 0
	cmdVideoStreamMode byte = 1

	// eARNETWORK_MANAGER_INTERNAL_BUFFER_ID
	ARNETWORK_MANAGER_INTERNAL_BUFFER_ID_PING byte = 0
	ARNETWORK_MANAGER_INTERNAL_BUFFER_ID_PONG byte = 1
)
