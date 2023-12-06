// type制御するために追加。本来ならばきっとemitはcomposableに入らないほうが良いのかも
type Emit = (event: (string), ...args: any[]) => void

const DEFAULT_PROCESS_COLOR = "rgb(66, 165, 246)"