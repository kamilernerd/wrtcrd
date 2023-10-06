type DatagramPayload = { Event: string, Value: any }

const KeyDownListener = (ev: KeyboardEvent, pc: RTCDataChannel) => {
  ev.preventDefault();
  ev.stopPropagation();
  if (pc.readyState !== "open") {
    return false
  }
  pc.send(Datagram({
    Event: "keyboard",
    Value: {
      State: 1,
      Key: ev.code.toLowerCase(),
      Shift: ev.shiftKey,
      Meta: ev.metaKey,
      Ctrl: ev.ctrlKey,
      Alt: ev.altKey
    }
  }))
  return
}

const KeyUpListener = (ev: KeyboardEvent, pc: RTCDataChannel) => {
  ev.preventDefault();
  ev.stopPropagation();
  if (pc.readyState !== "open") {
    return false
  }
  pc.send(Datagram({
    Event: "keyboard",
    Value: {
      State: 0,
      Key: ev.code.toLowerCase(),
      Shift: ev.shiftKey,
      Meta: ev.metaKey,
      Ctrl: ev.ctrlKey,
      Alt: ev.altKey
    }
  }))
  return
}

const Datagram = (data: DatagramPayload) => JSON.stringify(data)

export {
  KeyUpListener,
  KeyDownListener
}
