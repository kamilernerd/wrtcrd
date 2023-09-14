type SocketMessage = {
  Event: string
  Value: any
}

declare global {
  interface Navigator {
    app: any;
  }
}