// 这里的代码是从github上的一个项目拿来的，很抱歉我忘记是哪个项目了。


// 角度转弧度
// Math.PI = 180 度
export function angleToRadian(angle: number) {
    return angle * Math.PI / 180
}

/**
 * 计算根据圆心旋转后的点的坐标
 * param   {Object}  point  旋转前的点坐标
 * param   {Object}  center 旋转中心
 * param   {Number}  rotate 旋转的角度
 * return  {Object}         旋转后的坐标
 * https://www.zhihu.com/question/67425734/answer/252724399 旋转矩阵公式
 */
export function rotatedPointCoordinate(point: any, center: any, rotate: number) {
    /**
     * 旋转公式：
     *  点a(x, y)
     *  旋转中心c(x, y)
     *  旋转后点n(x, y)
     *  旋转角度θ                tan ??
     * nx = cosθ * (ax - cx) - sinθ * (ay - cy) + cx
     * ny = sinθ * (ax - cx) + cosθ * (ay - cy) + cy
     */

    return {
        x: (point.x - center.x) * Math.cos(angleToRadian(rotate)) - (point.y - center.y) * Math.sin(angleToRadian(rotate)) + center.x,
        y: (point.x - center.x) * Math.sin(angleToRadian(rotate)) + (point.y - center.y) * Math.cos(angleToRadian(rotate)) + center.y,
    }
}

// 求两点之间的中点坐标
export function getCenterPoint(p1: any, p2: any) {
    return {
        x: p1.x + ((p2.x - p1.x) / 2),
        y: p1.y + ((p2.y - p1.y) / 2),
    }
}

const funcs = {
    lt: calculateLeftTop,
    t: calculateTop,
    rt: calculateRightTop,
    r: calculateRight,
    rb: calculateRightBottom,
    b: calculateBottom,
    lb: calculateLeftBottom,
    l: calculateLeft,
}

function calculateLeftTop(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint } = pointInfo
    let newCenterPoint = getCenterPoint(curPositon, symmetricPoint)
    let newTopLeftPoint = rotatedPointCoordinate(curPositon, newCenterPoint, -style.rotate)
    let newBottomRightPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)

    let newWidth = newBottomRightPoint.x - newTopLeftPoint.x
    let newHeight = newBottomRightPoint.y - newTopLeftPoint.y

    if (needLockProportion) {
        if (newWidth / newHeight > proportion) {
            newTopLeftPoint.x += Math.abs(newWidth - newHeight * proportion)
            newWidth = newHeight * proportion
        } else {
            newTopLeftPoint.y += Math.abs(newHeight - newWidth / proportion)
            newHeight = newWidth / proportion
        }

        // 由于现在求的未旋转前的坐标是以没按比例缩减宽高前的坐标来计算的
        // 所以缩减宽高后，需要按照原来的中心点旋转回去，获得缩减宽高并旋转后对应的坐标
        // 然后以这个坐标和对称点获得新的中心点，并重新计算未旋转前的坐标
        const rotatedTopLeftPoint = rotatedPointCoordinate(newTopLeftPoint, newCenterPoint, style.rotate)
        newCenterPoint = getCenterPoint(rotatedTopLeftPoint, symmetricPoint)
        newTopLeftPoint = rotatedPointCoordinate(rotatedTopLeftPoint, newCenterPoint, -style.rotate)
        newBottomRightPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)

        newWidth = newBottomRightPoint.x - newTopLeftPoint.x
        newHeight = newBottomRightPoint.y - newTopLeftPoint.y
    }

    if (newWidth > 0 && newHeight > 0) {
        style.width = Math.round(newWidth)
        style.height = Math.round(newHeight)
        style.left = Math.round(newTopLeftPoint.x)
        style.top = Math.round(newTopLeftPoint.y)
    }
}

function calculateRightTop(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint } = pointInfo
    let newCenterPoint = getCenterPoint(curPositon, symmetricPoint)
    let newTopRightPoint = rotatedPointCoordinate(curPositon, newCenterPoint, -style.rotate)
    let newBottomLeftPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)

    let newWidth = newTopRightPoint.x - newBottomLeftPoint.x
    let newHeight = newBottomLeftPoint.y - newTopRightPoint.y

    if (needLockProportion) {
        if (newWidth / newHeight > proportion) {
            newTopRightPoint.x -= Math.abs(newWidth - newHeight * proportion)
            newWidth = newHeight * proportion
        } else {
            newTopRightPoint.y += Math.abs(newHeight - newWidth / proportion)
            newHeight = newWidth / proportion
        }

        const rotatedTopRightPoint = rotatedPointCoordinate(newTopRightPoint, newCenterPoint, style.rotate)
        newCenterPoint = getCenterPoint(rotatedTopRightPoint, symmetricPoint)
        newTopRightPoint = rotatedPointCoordinate(rotatedTopRightPoint, newCenterPoint, -style.rotate)
        newBottomLeftPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)

        newWidth = newTopRightPoint.x - newBottomLeftPoint.x
        newHeight = newBottomLeftPoint.y - newTopRightPoint.y
    }

    if (newWidth > 0 && newHeight > 0) {
        style.width = Math.round(newWidth)
        style.height = Math.round(newHeight)
        style.left = Math.round(newBottomLeftPoint.x)
        style.top = Math.round(newTopRightPoint.y)
    }
}

function calculateRightBottom(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint } = pointInfo
    let newCenterPoint = getCenterPoint(curPositon, symmetricPoint)
    let newTopLeftPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)
    let newBottomRightPoint = rotatedPointCoordinate(curPositon, newCenterPoint, -style.rotate)

    let newWidth = newBottomRightPoint.x - newTopLeftPoint.x
    let newHeight = newBottomRightPoint.y - newTopLeftPoint.y

    if (needLockProportion) {
        if (newWidth / newHeight > proportion) {
            newBottomRightPoint.x -= Math.abs(newWidth - newHeight * proportion)
            newWidth = newHeight * proportion
        } else {
            newBottomRightPoint.y -= Math.abs(newHeight - newWidth / proportion)
            newHeight = newWidth / proportion
        }

        const rotatedBottomRightPoint = rotatedPointCoordinate(newBottomRightPoint, newCenterPoint, style.rotate)
        newCenterPoint = getCenterPoint(rotatedBottomRightPoint, symmetricPoint)
        newTopLeftPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)
        newBottomRightPoint = rotatedPointCoordinate(rotatedBottomRightPoint, newCenterPoint, -style.rotate)

        newWidth = newBottomRightPoint.x - newTopLeftPoint.x
        newHeight = newBottomRightPoint.y - newTopLeftPoint.y
    }

    if (newWidth > 0 && newHeight > 0) {
        style.width = Math.round(newWidth)
        style.height = Math.round(newHeight)
        style.left = Math.round(newTopLeftPoint.x)
        style.top = Math.round(newTopLeftPoint.y)
    }
}

function calculateLeftBottom(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint } = pointInfo
    let newCenterPoint = getCenterPoint(curPositon, symmetricPoint)
    let newTopRightPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)
    let newBottomLeftPoint = rotatedPointCoordinate(curPositon, newCenterPoint, -style.rotate)

    let newWidth = newTopRightPoint.x - newBottomLeftPoint.x
    let newHeight = newBottomLeftPoint.y - newTopRightPoint.y

    if (needLockProportion) {
        if (newWidth / newHeight > proportion) {
            newBottomLeftPoint.x += Math.abs(newWidth - newHeight * proportion)
            newWidth = newHeight * proportion
        } else {
            newBottomLeftPoint.y -= Math.abs(newHeight - newWidth / proportion)
            newHeight = newWidth / proportion
        }

        const rotatedBottomLeftPoint = rotatedPointCoordinate(newBottomLeftPoint, newCenterPoint, style.rotate)
        newCenterPoint = getCenterPoint(rotatedBottomLeftPoint, symmetricPoint)
        newTopRightPoint = rotatedPointCoordinate(symmetricPoint, newCenterPoint, -style.rotate)
        newBottomLeftPoint = rotatedPointCoordinate(rotatedBottomLeftPoint, newCenterPoint, -style.rotate)

        newWidth = newTopRightPoint.x - newBottomLeftPoint.x
        newHeight = newBottomLeftPoint.y - newTopRightPoint.y
    }

    if (newWidth > 0 && newHeight > 0) {
        style.width = Math.round(newWidth)
        style.height = Math.round(newHeight)
        style.left = Math.round(newBottomLeftPoint.x)
        style.top = Math.round(newTopRightPoint.y)
    }
}

function calculateTop(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint, curPoint } = pointInfo
    // 由于用户拉伸时是以任意角度拉伸的，所以在求得旋转前的坐标时，只取 y 坐标（这里的 x 坐标可能是任意值），x 坐标用 curPoint 的。
    // 这个中心点（第二个参数）用 curPoint, center, symmetricPoint 都可以，只要他们在一条直线上就行
    const rotatedcurPositon = rotatedPointCoordinate(curPositon, curPoint, -style.rotate)

    // 算出旋转前 y 坐标，再用 curPoint 的 x 坐标，重新计算它们旋转后对应的坐标
    const rotatedTopMiddlePoint = rotatedPointCoordinate({
        x: curPoint.x,
        y: rotatedcurPositon.y,
    }, curPoint, style.rotate)

    // 用旋转后的坐标和对称点算出新的高度（勾股定理）
    const newHeight = Math.sqrt((rotatedTopMiddlePoint.x - symmetricPoint.x) ** 2 + (rotatedTopMiddlePoint.y - symmetricPoint.y) ** 2)

    const newCenter = {
        x: rotatedTopMiddlePoint.x - (rotatedTopMiddlePoint.x - symmetricPoint.x) / 2,
        y: rotatedTopMiddlePoint.y + (symmetricPoint.y - rotatedTopMiddlePoint.y) / 2,
    }

    let width = style.width
    // 因为调整的是高度 所以只需根据锁定的比例调整宽度即可
    if (needLockProportion) {
        width = newHeight * proportion
    }

    style.width = width
    style.height = Math.round(newHeight)
    style.top = Math.round(newCenter.y - (newHeight / 2))
    style.left = Math.round(newCenter.x - (style.width / 2))
}

function calculateRight(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint, curPoint } = pointInfo
    const rotatedcurPositon = rotatedPointCoordinate(curPositon, curPoint, -style.rotate)
    const rotatedRightMiddlePoint = rotatedPointCoordinate({
        x: rotatedcurPositon.x,
        y: curPoint.y,
    }, curPoint, style.rotate)

    const newWidth = Math.sqrt((rotatedRightMiddlePoint.x - symmetricPoint.x) ** 2 + (rotatedRightMiddlePoint.y - symmetricPoint.y) ** 2)

    const newCenter = {
        x: rotatedRightMiddlePoint.x - (rotatedRightMiddlePoint.x - symmetricPoint.x) / 2,
        y: rotatedRightMiddlePoint.y + (symmetricPoint.y - rotatedRightMiddlePoint.y) / 2,
    }

    let height = style.height
    // 因为调整的是宽度 所以只需根据锁定的比例调整高度即可
    if (needLockProportion) {
        height = newWidth / proportion
    }

    style.height = height
    style.width = Math.round(newWidth)
    style.top = Math.round(newCenter.y - (style.height / 2))
    style.left = Math.round(newCenter.x - (newWidth / 2))
}

function calculateBottom(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint, curPoint } = pointInfo
    const rotatedcurPositon = rotatedPointCoordinate(curPositon, curPoint, -style.rotate)
    const rotatedBottomMiddlePoint = rotatedPointCoordinate({
        x: curPoint.x,
        y: rotatedcurPositon.y,
    }, curPoint, style.rotate)

    const newHeight = Math.sqrt((rotatedBottomMiddlePoint.x - symmetricPoint.x) ** 2 + (rotatedBottomMiddlePoint.y - symmetricPoint.y) ** 2)

    const newCenter = {
        x: rotatedBottomMiddlePoint.x - (rotatedBottomMiddlePoint.x - symmetricPoint.x) / 2,
        y: rotatedBottomMiddlePoint.y + (symmetricPoint.y - rotatedBottomMiddlePoint.y) / 2,
    }

    let width = style.width
    // 因为调整的是高度 所以只需根据锁定的比例调整宽度即可
    if (needLockProportion) {
        width = newHeight * proportion
    }

    style.width = width
    style.height = Math.round(newHeight)
    style.top = Math.round(newCenter.y - (newHeight / 2))
    style.left = Math.round(newCenter.x - (style.width / 2))
}

function calculateLeft(style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    const { symmetricPoint, curPoint } = pointInfo
    const rotatedcurPositon = rotatedPointCoordinate(curPositon, curPoint, -style.rotate)
    const rotatedLeftMiddlePoint = rotatedPointCoordinate({
        x: rotatedcurPositon.x,
        y: curPoint.y,
    }, curPoint, style.rotate)

    const newWidth = Math.sqrt((rotatedLeftMiddlePoint.x - symmetricPoint.x) ** 2 + (rotatedLeftMiddlePoint.y - symmetricPoint.y) ** 2)

    const newCenter = {
        x: rotatedLeftMiddlePoint.x - (rotatedLeftMiddlePoint.x - symmetricPoint.x) / 2,
        y: rotatedLeftMiddlePoint.y + (symmetricPoint.y - rotatedLeftMiddlePoint.y) / 2,
    }

    let height = style.height
    if (needLockProportion) {
        height = newWidth / proportion
    }

    style.height = height
    style.width = Math.round(newWidth)
    style.top = Math.round(newCenter.y - (style.height / 2))
    style.left = Math.round(newCenter.x - (newWidth / 2))
}

export default function positonAndSize(name: string, style: any, curPositon: any, proportion: number, needLockProportion: boolean, pointInfo: any) {
    funcs[name as keyof typeof funcs](style, curPositon, proportion, needLockProportion, pointInfo)
}