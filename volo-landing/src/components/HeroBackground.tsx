import { useEffect, useRef } from "react";
import * as THREE from "three";

/**
 * WebGL hero background using Three.js — plasma wave interference shader.
 * Inspired by Raycast's implementation: a ShaderMaterial on a plane with
 * animated color mixing driven by time-based trigonometric functions.
 *
 * Adapted for Volo's sky-blue/teal/indigo palette.
 */

const VERTEX_SHADER = `
varying vec2 vUv;

void main() {
    vUv = uv;
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
}
`;

const FRAGMENT_SHADER = `
uniform float uTime;
uniform float uSpeed;
uniform float uSize;
uniform float uBrightness;
uniform vec3 uColor1;
uniform vec3 uColor2;
uniform vec3 uColor3;
uniform vec2 uResolution;

varying vec2 vUv;

#define PI 3.14159265359

void main() {
    vec2 uv = gl_FragCoord.xy / uResolution.xy;
    float aspect = uResolution.x / uResolution.y;
    vec2 coord = vec2(uv.x * aspect, uv.y);

    float t = uTime * uSpeed;

    // Diagonal axis — rotate coordinates ~40 degrees
    float angle = 0.7;
    float ca = cos(angle);
    float sa = sin(angle);
    vec2 rotated = vec2(
        coord.x * ca - coord.y * sa,
        coord.x * sa + coord.y * ca
    );

    // Create repeating bands along the rotated x-axis
    float bandFreq = uSize * 3.0;
    float bandPos = fract(rotated.x * bandFreq);

    // Shape each band as a cylinder cross-section
    float tubeWidth = 0.55;
    float dist = abs(bandPos - 0.5) / tubeWidth;
    float tube = 1.0 - smoothstep(0.0, 1.0, dist);

    // Specular highlight — sweeps along the tube
    float highlightPos = sin(t * 0.4 + rotated.x * bandFreq * PI * 0.3 + rotated.y * 2.0) * 0.5 + 0.5;
    float specular = pow(1.0 - abs(bandPos - highlightPos), 8.0) * tube;

    // Secondary moving reflection
    float highlight2 = sin(t * 0.25 + rotated.x * bandFreq * 0.7 - rotated.y * 1.5) * 0.5 + 0.5;
    float specular2 = pow(1.0 - abs(bandPos - highlight2), 12.0) * tube * 0.6;

    // Per-band hue shift
    float bandIndex = floor(rotated.x * bandFreq);
    float hueShift = sin(bandIndex * 1.7 + t * 0.1) * 0.5 + 0.5;

    // Base tube color
    vec3 tubeColor = mix(uColor3, uColor2, tube * 0.7);

    // Hue variation per band
    vec3 bandTint = mix(uColor2, uColor1, hueShift);
    tubeColor = mix(tubeColor, bandTint, tube * 0.4);

    // Specular highlights boosted by brightness
    vec3 specColor = mix(uColor1, vec3(1.0, 0.95, 0.9), 0.4);
    tubeColor += specColor * specular * (1.0 + uBrightness * 1.5);
    tubeColor += uColor1 * specular2 * (0.8 + uBrightness * 1.0);

    // Ambient brightness lift — illuminates the whole surface
    tubeColor += uColor2 * tube * uBrightness * 0.6;

    // Edge darkening
    float edgeDark = smoothstep(0.0, 0.3, tube);
    tubeColor *= edgeDark;

    // Gamma
    tubeColor = pow(tubeColor, vec3(1.3 - uBrightness * 0.2));

    gl_FragColor = vec4(tubeColor, 1.0);
}
`;

// Volo palette — metallic ridges with sky-blue/cyan specular highlights
const COLORS = {
  color1: new THREE.Color(0.10, 0.60, 1.0),   // Sky-blue highlight (specular peak)
  color2: new THREE.Color(0.03, 0.15, 0.35),   // Deep blue-steel midtone
  color3: new THREE.Color(0.005, 0.01, 0.03),  // Near-black gap/shadow
};

export function HeroBackground({ className = "" }: { className?: string }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const rendererRef = useRef<THREE.WebGLRenderer | null>(null);
  const frameRef = useRef<number>(0);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    // Setup
    const scene = new THREE.Scene();
    const camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0, 1);

    const renderer = new THREE.WebGLRenderer({
      antialias: false,
      alpha: true,
      powerPreference: "high-performance",
    });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.setSize(container.clientWidth, container.clientHeight);
    container.appendChild(renderer.domElement);
    rendererRef.current = renderer;

    // Shader material
    const uniforms = {
      uTime: { value: 0 },
      uSpeed: { value: 0.3 },
      uSize: { value: 2.5 },
      uBrightness: { value: 0 },
      uColor1: { value: COLORS.color1 },
      uColor2: { value: COLORS.color2 },
      uColor3: { value: COLORS.color3 },
      uResolution: { value: new THREE.Vector2(container.clientWidth, container.clientHeight) },
    };

    const material = new THREE.ShaderMaterial({
      vertexShader: VERTEX_SHADER,
      fragmentShader: FRAGMENT_SHADER,
      uniforms,
    });

    const geometry = new THREE.PlaneGeometry(2, 2);
    const mesh = new THREE.Mesh(geometry, material);
    scene.add(mesh);

    // Animation loop
    const clock = new THREE.Clock();

    function animate() {
      frameRef.current = requestAnimationFrame(animate);
      const elapsed = clock.getElapsedTime();
      uniforms.uTime.value = elapsed;

      // Brightness pulse — peaks every ~8 seconds, smooth sine curve
      // Goes from 0 (dark) to 1 (lit) and back
      const cycle = 8.0; // seconds per full cycle
      const raw = Math.sin(elapsed * (2 * Math.PI / cycle) - Math.PI / 2); // starts low
      // Shape it: mostly dark, brief bright peak
      const shaped = Math.max(0, raw);
      uniforms.uBrightness.value = shaped * shaped; // square for sharper peak

      renderer.render(scene, camera);
    }
    animate();

    // Resize handler
    function onResize() {
      if (!container || !renderer) return;
      const w = container.clientWidth;
      const h = container.clientHeight;
      renderer.setSize(w, h);
      uniforms.uResolution.value.set(w, h);
    }

    window.addEventListener("resize", onResize);

    return () => {
      cancelAnimationFrame(frameRef.current);
      window.removeEventListener("resize", onResize);
      renderer.dispose();
      geometry.dispose();
      material.dispose();
      if (container.contains(renderer.domElement)) {
        container.removeChild(renderer.domElement);
      }
    };
  }, []);

  return (
    <div
      ref={containerRef}
      className={`absolute inset-0 ${className}`}
      style={{ pointerEvents: "none" }}
    />
  );
}
