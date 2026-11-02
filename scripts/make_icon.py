#!/usr/bin/env python3
"""Draws the app icon: a server chassis with three rack units as a face,
framed by long, straight, black curtain hair parted in the middle.

It writes three files:
  assets/icon.png                the icon on its rounded green tile
  assets/android-background.png  Android's adaptive icon background layer
  assets/android-foreground.png  Android's adaptive icon foreground layer: the
                                 character alone, inside the 66dp safe zone of
                                 the 108dp canvas, so launchers can mask the
                                 icon to any shape without edges showing

Run from the repository root: python3 scripts/make_icon.py (needs Pillow)."""
from PIL import Image, ImageDraw, ImageFilter

S = 1024          # output size
K = 4             # supersampling factor
W = S * K


def P(x, y):
    """Converts unit coordinates to pixels."""
    return (x * W, y * W)


def lerp(a, b, t):
    return tuple(int(a[i] + (b[i] - a[i]) * t) for i in range(len(a)))


def bezier(p0, p1, p2, p3, n=60):
    """Points along a cubic Bézier curve, in unit coordinates."""
    pts = []
    for i in range(n + 1):
        t = i / n
        u = 1 - t
        x = u**3 * p0[0] + 3 * u * u * t * p1[0] + 3 * u * t * t * p2[0] + t**3 * p3[0]
        y = u**3 * p0[1] + 3 * u * u * t * p1[1] + 3 * u * t * t * p2[1] + t**3 * p3[1]
        pts.append(P(x, y))
    return pts


def path(*segments):
    pts = []
    for s in segments:
        pts.extend(bezier(*s))
    return pts


def mirror(segments):
    return [tuple((1 - x, y) for x, y in s) for s in segments]


# Background: a dark emerald tile.
bg = Image.new("RGBA", (W, W))
d = ImageDraw.Draw(bg)
for y in range(W):
    d.line([(0, y), (W, y)], fill=lerp((22, 58, 44), (5, 14, 11), y / W) + (255,))
glow = Image.new("L", (W, W), 0)
ImageDraw.Draw(glow).ellipse([*P(0.18, 0.12), *P(0.82, 0.72)], fill=110)
glow = glow.filter(ImageFilter.GaussianBlur(W * 0.08))
bg.paste(Image.new("RGBA", (W, W), (46, 150, 104, 255)), (0, 0), glow)
tile = Image.new("L", (W, W), 0)
ImageDraw.Draw(tile).rounded_rectangle([*P(0.07, 0.07), *P(0.93, 0.93)], radius=int(W * 0.22), fill=255)
background = bg

# The character is drawn on its own transparent layer, composed onto the
# tile at the end.
img = Image.new("RGBA", (W, W), (0, 0, 0, 0))

HAIR = (13, 13, 17, 255)
STRAND = (44, 46, 56)
SHEEN = (125, 135, 156)


def shadow_of(poly, dx=0.006, dy=0.012, blur=0.014, alpha=150):
    layer = Image.new("RGBA", (W, W), (0, 0, 0, 0))
    ImageDraw.Draw(layer).polygon([(x + dx * W, y + dy * W) for x, y in poly], fill=(0, 0, 0, alpha))
    return layer.filter(ImageFilter.GaussianBlur(W * blur))


# The outside of one side of the hair, from the parting over the crown
# and down to the shoulder.
outer = [
    ((0.500, 0.160), (0.400, 0.158), (0.262, 0.190), (0.236, 0.330)),
    ((0.236, 0.330), (0.222, 0.450), (0.222, 0.640), (0.212, 0.800)),
]
# The inside: from the shoulder up along the cheek, over the temple and
# back to the parting.
inner = [
    ((0.375, 0.800), (0.372, 0.650), (0.370, 0.480), (0.372, 0.360)),
    ((0.372, 0.360), (0.376, 0.290), (0.420, 0.225), (0.500, 0.205)),
]
# Ragged ends between the outer and inner edges.
ends = [(0.212, 0.800), (0.222, 0.832), (0.248, 0.842), (0.272, 0.830), (0.300, 0.852),
        (0.330, 0.838), (0.356, 0.828), (0.375, 0.800)]


def side(flip):
    def f(pts):
        return [((W - x) if flip else x, y) for x, y in pts]
    poly = path(*outer) + [P(x, y) for x, y in ends] + path(*inner)
    return f(poly)


# Back hair: the crown behind the head, so nothing shows between the
# curtains above the face.
crown = path(*outer) + [P(0.215, 0.70), P(0.785, 0.70)] + list(reversed([(W - x, y) for x, y in path(*outer)]))
img = Image.alpha_composite(img, shadow_of(crown, dy=0.02, blur=0.025))
d = ImageDraw.Draw(img)
d.polygon(crown, fill=HAIR)

# The face: a graphite chassis holding three rack units.
fx0, fx1, fy0, fy1 = 0.330, 0.670, 0.250, 0.780
face = Image.new("RGBA", (W, W), (0, 0, 0, 0))
fd = ImageDraw.Draw(face)
for y in range(int(fy0 * W), int(fy1 * W)):
    fd.line([(fx0 * W, y), (fx1 * W, y)], fill=lerp((62, 70, 74), (34, 40, 43), (y / W - fy0) / (fy1 - fy0)) + (255,))
m = Image.new("L", (W, W), 0)
ImageDraw.Draw(m).rounded_rectangle([*P(fx0, fy0), *P(fx1, fy1)], radius=int(W * 0.05), fill=255)
img = Image.alpha_composite(img, shadow_of([P(fx0, fy0), P(fx1, fy0), P(fx1, fy1), P(fx0, fy1)], dy=0.015, blur=0.02))
img.paste(face, (0, 0), m)
d = ImageDraw.Draw(img)

x0, x1 = 0.372, 0.628
silver_top, silver_bottom = (220, 229, 225), (162, 178, 172)
for i in range(3):
    y0 = 0.330 + i * 0.140
    y1 = y0 + 0.108
    unit = Image.new("RGBA", (W, W), (0, 0, 0, 0))
    ud = ImageDraw.Draw(unit)
    for y in range(int(y0 * W), int(y1 * W)):
        ud.line([(x0 * W, y), (x1 * W, y)], fill=lerp(silver_top, silver_bottom, (y / W - y0) / (y1 - y0)) + (255,))
    um = Image.new("L", (W, W), 0)
    ImageDraw.Draw(um).rounded_rectangle([*P(x0, y0), *P(x1, y1)], radius=int(W * 0.018), fill=255)
    img.paste(unit, (0, 0), um)
    d = ImageDraw.Draw(img)
    cy = (y0 + y1) / 2
    d.rounded_rectangle([*P(x0 + 0.105, cy - 0.008), *P(x1 - 0.035, cy + 0.008)], radius=int(W * 0.008), fill=(18, 26, 23, 255))
    r = 0.014
    d.ellipse([*P(x0 + 0.050 - r, cy - r), *P(x0 + 0.050 + r, cy + r)], fill=(36, 214, 128, 255))

# Front curtains over the temples and cheeks.
for flip in (False, True):
    poly = side(flip)
    img = Image.alpha_composite(img, shadow_of(poly, dx=(-0.006 if flip else 0.006)))
    d = ImageDraw.Draw(img)
    d.polygon(poly, fill=HAIR)

# Strands: thin lines following the fall of the hair, and a greasy sheen
# where it curves over the crown.
lines = Image.new("RGBA", (W, W), (0, 0, 0, 0))
ld = ImageDraw.Draw(lines)
strands = [
    ((0.47, 0.172), (0.36, 0.180), (0.28, 0.230), (0.258, 0.400), (0.250, 0.640), (0.240, 0.815)),
    ((0.49, 0.190), (0.40, 0.205), (0.33, 0.260), (0.315, 0.420), (0.310, 0.650), (0.300, 0.840)),
    ((0.48, 0.200), (0.41, 0.225), (0.36, 0.290), (0.350, 0.440), (0.352, 0.640), (0.348, 0.830)),
    ((0.46, 0.180), (0.38, 0.195), (0.30, 0.250), (0.280, 0.420), (0.272, 0.660), (0.262, 0.835)),
]
for st in strands:
    for flip in (False, True):
        a, b, c, e, f, g = st
        pts = bezier(a, b, c, e, n=50) + bezier(e, f, f, g, n=50)
        if flip:
            pts = [(W - x, y) for x, y in pts]
        ld.line(pts, fill=STRAND + (255,), width=int(W * 0.004))
for flip in (False, True):
    pts = bezier((0.47, 0.168), (0.38, 0.170), (0.30, 0.200), (0.268, 0.300), n=80)
    if flip:
        pts = [(W - x, y) for x, y in pts]
    ld.line(pts, fill=SHEEN + (170,), width=int(W * 0.010))
lines = lines.filter(ImageFilter.GaussianBlur(W * 0.003))
img = Image.alpha_composite(img, lines)

# The parting.
d = ImageDraw.Draw(img)
d.line([P(0.5, 0.158), P(0.5, 0.215)], fill=(70, 74, 84, 255), width=int(W * 0.005))

character = img

icon = Image.new("RGBA", (W, W), (0, 0, 0, 0))
icon.paste(background, (0, 0), tile)
icon = Image.alpha_composite(icon, character)
icon.resize((S, S), Image.LANCZOS).save("assets/icon.png")

background.convert("RGB").resize((S, S), Image.LANCZOS).save("assets/android-background.png")

# The character spans about 0.62 x 0.72 of the canvas; scaled by 0.66 its
# corners stay inside the circle launchers show.
scale = 0.66
small = character.resize((int(W * scale), int(W * scale)), Image.LANCZOS)
foreground = Image.new("RGBA", (W, W), (0, 0, 0, 0))
offset = int(W * (1 - scale) / 2)
foreground.alpha_composite(small, (offset, offset - int(W * 0.01)))
foreground.resize((S, S), Image.LANCZOS).save("assets/android-foreground.png")
