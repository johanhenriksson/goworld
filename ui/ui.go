package ui;

import (
)

type Drawable interface {
    Draw()
}

type Rect struct {
    *Element
}

