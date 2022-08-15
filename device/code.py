from rainbowio import colorwheel
from adafruit_macropad import MacroPad
import usb_cdc
import time

macropad = MacroPad()

keys = macropad.keys

button = macropad.encoder_switch
pixels = macropad.pixels
mode = "0" # Mode 0 = bootstrap.
submode = "0" # Mode 0 = default
serial = usb_cdc.data
serialbuffer = ""
serialdata = ""
hbtime = 0
hbsendtime = 0
tones = []

inptypes = {"serial": 0, "key": 1, "rotary": 2, "rotary_btn": 3}

inpCnt = 0

last_position = 0
lastvlu = False

lastmsg = ""

def writeScreen(msg):
    global lastmsg
    if msg == lastmsg:
        return False
    txtlines = macropad.display_text(title="Macropad v0.0.1")
    txtlines[0].text = msg
    txtlines.show()
    lastmsg = msg

def sendSerial(inptype, data, increment=True):
    global inpCnt

    out = b"%d.%d.%s.%s.%s\n" %(inpCnt, inptypes[inptype], mode, submode, data)
    serial.write(out)
    if increment:
        inpCnt += 1

def _getSerialRaw():
    text = ""
    available = serial.in_waiting
    while available:
        raw = serial.read(available)
        text = raw.decode("utf-8")
        available = serial.in_waiting
    return text

def _getSerialData():
    global serialbuffer, serialdata
    serialbuffer += _getSerialRaw()
    if serialbuffer.endswith("\n"):
        # strip line end
        serialdata = serialbuffer[:-1]
        # clear buffer
        serialbuffer = ""

def setMode(newmode, newsub):
    global mode
    global submode
    mode = newmode

    data = "Mode: %s" %mode

    if submode != "" and submode:
        submode = newsub
        data += ", Submode: %s" %submode

    writeScreen(data)

def getSerial(now):
    global hbtime, serialdata, tones

    _getSerialData()

    if not serialdata:
        return False
    if "heartbeat" in serialdata:
        serialdata = ""
        hbtime = now
    elif "mode." in serialdata:
        data = serialdata.split(".")
        subm = ""
        if len(data) == 3:
            subm = data[2]
        setMode(data[1], subm)
    elif "color." in serialdata:
        colors = serialdata.split(".")[1:]
        if len(colors) == 12:
            data = []
            for item in colors:
                data.append(int(item))
            setColor(data)
    elif "brightness." in serialdata:
        i = 0
        for item in serialdata.split(".")[1:]:
            i = float(item.replace("-","."))
        pixels.brightness = i
    elif "tone." in serialdata:
        data = serialdata.split(".")[1:]
        if len(data) == 12:
            tones = []
            for item in data:
                tones.append(int(item))


def setColor(colorArr):
    for i in range(len(colorArr)):
        pixels[i] = colorwheel(colorArr[i])

def readInput(inptype, data):
    if inptype == "key" and tones:
        macropad.play_tone(tones[data], 0.1)
    sendSerial(inptype, data)

def sendHeartbeat(now):
    global hbsendtime
    if hbsendtime == 0 or now - hbsendtime > 1: #If it's been more than a second since the last heartbeat
        sendSerial("serial", "heartbeat")
        if mode == "0": # Also send a bootstrap message if we're in that mode
            sendSerial("serial", "bootstrap")
        hbsendtime = now

def checkheartBeat(now): # Should get heartbeat every second from server, driver for connectivity to driver
    global serialdata, hbtime, mode
    if hbtime == 0:
        return True
    if now - hbtime > 3: # If more than 5 seconds since heartbeat put back in bootstrap
        mode = "0"
        setDefaultMacroPosition()
        writeScreen("Lost connection to PC")

def getInput():
    global last_position
    global lastvlu

    if mode == "0":

        return False

    button = macropad.encoder_switch
    if lastvlu != button:
        lastvlu = button
        readInput("rotary_btn", int(button))
        if not button:
            pixels.brightness = 1.0
        else:
            pixels.brightness = 0.2

    
    position = macropad.encoder
    if last_position is None or position != last_position:
        readInput("rotary", position)
    last_position = position

    event = keys.events.get()
    if event:
        if event.pressed:
            readInput("key", event.key_number)

def setDefaultMacroPosition():
    pixels.fill([255,255,0])
    pixels.brightness = 0.2

def main():

    setDefaultMacroPosition()
    writeScreen("Connecting to PC")

    while True:
        now = time.monotonic()
        sendHeartbeat(now)

        if mode == "0": #Clear our events in bootstrap
            keys.events.clear() 
        else:
            checkheartBeat(now) #Check our heartbeat
            
        getSerial(now)
        getInput()

if __name__ == "__main__":
    main()