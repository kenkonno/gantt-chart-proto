// type制御するために追加。本来ならばきっとemitはcomposableに入らないほうが良いのかも
type Emit = (event: (string), ...args: any[]) => void