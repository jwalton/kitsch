import React from "react";

type Color = [number, number, number];

interface ForegroundGradientProps {
  gradientStart: string;
  gradientEnd: string;
  children: string;
}

function hexToColors(hex: string): Color {
  const hexes = hex.match(/.{1,2}/g);
  const [red, green, blue] = [
    parseInt(hexes[0], 16),
    parseInt(hexes[1], 16),
    parseInt(hexes[2], 16),
  ];
  return [red, green, blue];
}

function lerp(c1: number, c2: number, s: number): number {
  return Math.floor(c1 + s * (c2 - c1));
}

function lerpColor(startColor: Color, endColor: Color, s: number): Color {
  return [
    lerp(startColor[0], endColor[0], s),
    lerp(startColor[1], endColor[1], s),
    lerp(startColor[2], endColor[2], s),
  ];
}

const ForegroundGradient: React.FC<ForegroundGradientProps> = (props) => {
  return <div>Hello World</div>;
};

export default React.memo(ForegroundGradient);
