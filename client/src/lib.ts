type DatagramPayload = { Event: string, Value: any }

const KeyDownListener = (ev: KeyboardEvent, pc: RTCDataChannel) => {
  ev.preventDefault();
  ev.stopPropagation();
  ev.stopImmediatePropagation()

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
  ev.stopImmediatePropagation()

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

const MouseEventListener = (ev: MouseEvent, pc: RTCDataChannel) => {
  ev.preventDefault()
  ev.stopPropagation()
  ev.stopImmediatePropagation()

  if (pc.readyState !== "open") {
    return false
  }

  pc.send(Datagram({
    Event: "move",
    Value: btoa(`${ev.currentTarget.id};${ev.pageX / ev.currentTarget!.clientWidth};${ev.offsetY / ev.currentTarget!.clientHeight}`)
  }))
}

const MouseDownListener = (ev: MouseEvent, pc: RTCDataChannel) => {
  ev.preventDefault()
  ev.stopPropagation()
  ev.stopImmediatePropagation()

  if (pc.readyState !== "open") {
    return false
  }

  pc.send(Datagram({
    Event: "action",
    Value: ev.button,
  }))
}

const MouseUpListener = (ev: MouseEvent, pc: RTCDataChannel) => {
  ev.preventDefault()
  ev.stopPropagation()
  ev.stopImmediatePropagation()

  if (pc.readyState !== "open") {
    return false
  }

  pc.send(Datagram({
    Event: "action",
    Value: ev.button,
  }))
}

const MouseScrollListener = (ev: WheelEvent, pc: RTCDataChannel) => {
  ev.preventDefault()
  ev.stopPropagation()
  ev.stopImmediatePropagation()

  if (pc.readyState !== "open") {
    return false
  }

  pc.send(Datagram({
    Event: "action",
    Value: String(`${ev.deltaX};${ev.deltaY}`) //x,y
  }))
}
const Datagram = (data: DatagramPayload) => JSON.stringify(data)

export {
  KeyUpListener,
  KeyDownListener,
  MouseEventListener,
  MouseDownListener,
  MouseUpListener,
  MouseScrollListener
}
