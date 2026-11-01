#!/usr/bin/env python3
"""Draws assets/icon.png: an emerald serpent coiled into an S around a server rack.
Run from the repository root: python3 scripts/make_icon.py (needs Pillow)."""
import math
from PIL import Image, ImageDraw, ImageFilter
S=1024; K=4; W=S*K
def lerp(a,b,t): return tuple(int(a[i]+(b[i]-a[i])*t) for i in range(len(a)))

# Background: deep emerald-black rounded square with a soft vignette.
bg=Image.new("RGBA",(W,W))
d=ImageDraw.Draw(bg)
top=(18,46,36); bot=(4,10,8)
for y in range(W):
    d.line([(0,y),(W,y)],fill=lerp(top,bot,y/W)+(255,))
glow=Image.new("L",(W,W),0)
ImageDraw.Draw(glow).ellipse([W*0.15,W*0.05,W*0.85,W*0.7],fill=90)
glow=glow.filter(ImageFilter.GaussianBlur(W*0.08))
bg.paste(Image.new("RGBA",(W,W),(40,140,95,255)),(0,0),glow)
mask=Image.new("L",(W,W),0)
m=int(W*0.07)
ImageDraw.Draw(mask).rounded_rectangle([m,m,W-m,W-m],radius=int(W*0.22),fill=255)
img=Image.new("RGBA",(W,W),(0,0,0,0)); img.paste(bg,(0,0),mask)
d=ImageDraw.Draw(img)
# Thin silver rim.

# Server rack: three silver units.
x0,x1=int(W*0.33),int(W*0.67)
silver=(176,192,186,255); slot=(10,22,18,255)
for i in range(3):
    y0=int(W*(0.33+i*0.13)); y1=y0+int(W*0.095)
    d.rounded_rectangle([x0,y0,x1,y1],radius=int(W*0.025),fill=silver)
    cy=(y0+y1)//2
    d.rounded_rectangle([x0+int(W*0.17),cy-int(W*0.009),x1-int(W*0.05),cy+int(W*0.009)],radius=int(W*0.009),fill=slot)
    r=int(W*0.016); cx=x0+int(W*0.065)
    d.ellipse([cx-r,cy-r,cx+r,cy+r],fill=(30,200,120,255))

# Serpent: an S-curve coiling around the rack, tapering to the tail.
def S_curve(t):
    # t in [0,1]: head top-right -> tail bottom-left
    a=math.pi*(1.15-2.3*t)
    x=0.5+0.27*math.cos(a)*(1 if t<0.5 else -1)
    y=0.20+0.62*t
    return x,y
pts=[]
N=1000
r=0.150
for i in range(N+1):
    t=i/N
    if t<0.5:
        th=math.radians(-25-245*(t/0.5)); cx,cy=0.5,0.35
    else:
        th=math.radians(-90+250*((t-0.5)/0.5)); cx,cy=0.5,0.35+2*r
    pts.append(((cx+r*math.cos(th))*W,(cy+r*math.sin(th))*W))
emer_hi=(80,230,150); emer_lo=(10,120,70)
shadow=Image.new("RGBA",(W,W),(0,0,0,0)); sd=ImageDraw.Draw(shadow)
body=Image.new("RGBA",(W,W),(0,0,0,0)); bd=ImageDraw.Draw(body)
for i,(x,y) in enumerate(pts):
    t=i/N
    w=W*(0.050*(1-t)**0.6+0.006)
    sd.ellipse([x-w+W*0.01,y-w+W*0.014,x+w+W*0.01,y+w+W*0.014],fill=(0,0,0,150))
    c=lerp(emer_hi,emer_lo,t)
    bd.ellipse([x-w,y-w,x+w,y+w],fill=c+(255,))
shadow=shadow.filter(ImageFilter.GaussianBlur(W*0.012))
img=Image.alpha_composite(img,shadow); img=Image.alpha_composite(img,body)
d=ImageDraw.Draw(img)
# A light stripe along the back gives the body some roundness.
stripe=Image.new("RGBA",(W,W),(0,0,0,0)); d=ImageDraw.Draw(stripe)
for i,(x,y) in enumerate(pts[:-120]):
    t=i/N
    w=W*(0.050*(1-t)**0.6+0.006)
    nx,ny=pts[i+3]
    ang=math.atan2(ny-y,nx-x)-math.pi/2
    ox,oy=x+math.cos(ang)*w*0.38,y+math.sin(ang)*w*0.38
    hw2=w*0.14
    d.ellipse([ox-hw2,oy-hw2,ox+hw2,oy+hw2],fill=(190,255,215,255))
stripe.putalpha(stripe.getchannel("A").point(lambda a: a*0.32))
img=Image.alpha_composite(img,stripe)
d=ImageDraw.Draw(img)
# Head at the start of the curve.
hx,hy=pts[0]; nx,ny=pts[12]
ang=math.atan2(hy-ny,hx-nx)
hw,hl=W*0.058,W*0.085
def rot(px,py): return (hx+px*math.cos(ang)-py*math.sin(ang), hy+px*math.sin(ang)+py*math.cos(ang))
hcx,hcy=rot(hl*0.35,0)
hd=Image.new("RGBA",(W,W),(0,0,0,0)); hdd=ImageDraw.Draw(hd)
hdd.ellipse([W/2-hl*0.75,W/2-hw,W/2+hl*0.75,W/2+hw],fill=emer_hi+(255,))
hd=hd.rotate(-math.degrees(ang),center=(W/2,W/2),resample=Image.BICUBIC)

img=Image.alpha_composite(img,hd.transform((W,W),Image.AFFINE,(1,0,W/2-hcx,0,1,W/2-hcy)))
d=ImageDraw.Draw(img)
ex,ey=rot(hl*0.55,-hw*0.42); er=W*0.014
d.ellipse([ex-er,ey-er,ex+er,ey+er],fill=(250,214,90,255))
d.ellipse([ex-er*0.35,ey-er,ex+er*0.35,ey+er],fill=(10,14,12,255))
# Forked tongue.
tx,ty=rot(hl*0.95,0); t1=rot(hl*1.4,-hw*0.2); t2=rot(hl*1.4,hw*0.2); tm=rot(hl*1.2,0)
d.line([ (tx,ty),tm],fill=(220,60,70,255),width=int(W*0.007))
d.line([tm,t1],fill=(220,60,70,255),width=int(W*0.006)); d.line([tm,t2],fill=(220,60,70,255),width=int(W*0.006))

img=img.resize((S,S),Image.LANCZOS)
img.save("assets/icon.png")
