#iChannel0 "file://textures/diffuse.jpg"    // Diffuse texture
#iChannel1 "file://textures/normal.jpg"     // Normal map
#iChannel2 "file://textures/height.jpg"     // Height map for bump mapping
#iChannel3 "file://Skyboxes/left.jpg"        // Background skybox

// Adjustable Parameters
const float BUMP_STRENGTH = 0.15;       // Strength of bump mapping (0.0 to 1.0)
const float NORMAL_STRENGTH = 1.2;        // Strength of normal mapping (0.0 to 2.0)
const float SPECULAR_POWER = 24.0;        // Specular highlight sharpness
const float AMBIENT_STRENGTH = 0.25;      // Ambient light intensity
const float DIFFUSE_STRENGTH = 1.2;         // Control diffuse light intensity
const float SPECULAR_INTENSITY = 0.8;       // Control specular intensity
const float RIM_STRENGTH = 0.15;            // Add rim lighting for edge highlights

// Animation Parameters
const vec3 ROTATION_SPEEDS = vec3(0.2, 0.3, 0.1);  // Rotation speed for each axis
const float LIGHT_ORBIT_SPEED = 0.5;       // Speed of light movement

// Object Parameters
const vec3 BOX_SIZE = vec3(1.0, 1.0, 1.0); // Box dimensions
const vec3 LIGHT_COLOR = vec3(1.0, 0.95, 0.8);

// Utility functions for rotation
mat3 rotateX(float angle) {
    float s = sin(angle);
    float c = cos(angle);
    return mat3(
        1, 0, 0,
        0, c, -s,
        0, s, c
    );
}

mat3 rotateY(float angle) {
    float s = sin(angle);
    float c = cos(angle);
    return mat3(
        c, 0, s,
        0, 1, 0,
        -s, 0, c
    );
}

mat3 rotateZ(float angle) {
    float s = sin(angle);
    float c = cos(angle);
    return mat3(
        c, -s, 0,
        s, c, 0,
        0, 0, 1
    );
}

// Box SDF
float sdBox(vec3 p, vec3 b) {
    vec3 q = abs(p) - b;
    return length(max(q, 0.0)) + min(max(q.x, max(q.y, q.z)), 0.0);
}

// Scene description
float scene(vec3 p) {
    // Apply rotation animation (same rotation used for UV mapping)
    p = rotateY(iTime * ROTATION_SPEEDS.x) * 
        rotateX(iTime * ROTATION_SPEEDS.y) * 
        rotateZ(iTime * ROTATION_SPEEDS.z) * p;
    
    return sdBox(p, BOX_SIZE);
}

// Calculate normal without bump/normal mapping
vec3 calcBaseNormal(vec3 p) {
    const float eps = 0.001;
    const vec2 h = vec2(eps, 0);
    return normalize(vec3(
        scene(p + h.xyy) - scene(p - h.xyy),
        scene(p + h.yxy) - scene(p - h.yxy),
        scene(p + h.yyx) - scene(p - h.yyx)
    ));
}

// --- FIX START ---
// New getBoxUV function that computes UV coordinates in object space.
// It applies the same rotation used in scene() so that the texture sticks to the cube.
vec2 getBoxUV(vec3 p) {
    // Compute the same rotation matrix as in scene()
    mat3 rot = rotateY(iTime * ROTATION_SPEEDS.x) *
               rotateX(iTime * ROTATION_SPEEDS.y) *
               rotateZ(iTime * ROTATION_SPEEDS.z);
    // Transform world-space point to object space
    vec3 localP = rot * p;
    
    vec2 uv;
    vec3 absP = abs(localP);
    if (absP.x > absP.y && absP.x > absP.z) {
        uv = localP.zy;
    } else if (absP.y > absP.z) {
        uv = localP.xz;
    } else {
        uv = localP.xy;
    }
    // Map from [-1,1] to [0,1]
    uv = uv * 0.5 + 0.5;
    return uv;
}
// --- FIX END ---

// Sample normal map and convert to world space
vec3 getNormalFromMap(vec2 uv, mat3 TBN) {
    vec3 normal = texture(iChannel1, uv).xyz * 2.0 - 1.0;
    normal = normalize(normal * vec3(NORMAL_STRENGTH, NORMAL_STRENGTH, 1.0));
    return normalize(TBN * normal);
}

// Calculate bump mapped normal
vec3 calcBumpedNormal(vec3 p, vec3 normal) {
    // --- FIX: Use getBoxUV(p) so UVs are computed in object space ---
    vec2 uv = getBoxUV(p);
    
    // Calculate tangent space based on the world-space normal
    vec3 tangent = normalize(cross(normal, vec3(0.0, 1.0, 0.0)));
    vec3 bitangent = normalize(cross(normal, tangent));
    mat3 TBN = mat3(tangent, bitangent, normal);
    
    // Sample height map for bump mapping
    float height = texture(iChannel2, uv).r;
    
    // Calculate bump mapped normal using finite differences
    const float bumpScale = BUMP_STRENGTH;
    const float h = 0.01;
    
    float hL = texture(iChannel2, uv - vec2(h, 0.0)).r;
    float hR = texture(iChannel2, uv + vec2(h, 0.0)).r;
    float hD = texture(iChannel2, uv - vec2(0.0, h)).r;
    float hU = texture(iChannel2, uv + vec2(0.0, h)).r;
    
    vec3 bumpNormal = normalize(vec3(
        (hL - hR) * bumpScale,
        (hD - hU) * bumpScale,
        2.0
    ));
    
    // Combine bump mapping with normal mapping
    vec3 normalMapNormal = getNormalFromMap(uv, TBN);
    return normalize(mix(normalMapNormal, bumpNormal, 0.5));
}

// Ray marching
float rayMarch(vec3 ro, vec3 rd) {
    float t = 0.0;
    for (int i = 0; i < 100; i++) {
        vec3 p = ro + rd * t;
        float d = scene(p);
        if (d < 0.001 || t > 20.0) break;
        t += d;
    }
    return t;
}

void mainImage(out vec4 fragColor, in vec2 fragCoord) {
    vec2 uv = fragCoord / iResolution.xy;
    vec3 ro = vec3(0.0, 0.0, -4.0);
    vec3 rd = normalize(vec3(((fragCoord - 0.5 * iResolution.xy) / iResolution.y).xy, 1.0));
    
    float t = rayMarch(ro, rd);
    
    if (t < 20.0) {
        vec3 p = ro + rd * t;
        vec3 baseNormal = calcBaseNormal(p);
        vec3 normal = calcBumpedNormal(p, baseNormal);
        
        // Dynamic light position
        vec3 lightPos = vec3(
            3.0 * cos(iTime * LIGHT_ORBIT_SPEED),
            2.0,
            3.0 * sin(iTime * LIGHT_ORBIT_SPEED)
        );
        
        vec3 lightDir = normalize(lightPos - p);
        vec3 viewDir = normalize(ro - p);
        vec3 halfDir = normalize(lightDir + viewDir);
        
        // --- FIX: Use getBoxUV(p) here so that the diffuse texture sticks to the cube face ---
        vec2 boxUV = getBoxUV(p);
        vec3 albedo = texture(iChannel0, boxUV).rgb;
        
        // Calculate lighting components
        float diff = max(dot(normal, lightDir), 0.0);
        float spec = pow(max(dot(normal, halfDir), 0.0), SPECULAR_POWER);
        
        // Add rim lighting (edge highlight)
        float rim = 1.0 - max(dot(viewDir, normal), 0.0);
        rim = pow(rim, 4.0) * RIM_STRENGTH;
        
        // Add fresnel effect to specular
        float fresnel = pow(1.0 - max(dot(normal, viewDir), 0.0), 5.0);
        
        // Combine lighting components with new parameters
        vec3 ambient = AMBIENT_STRENGTH * albedo;
        vec3 diffuse = DIFFUSE_STRENGTH * diff * albedo * LIGHT_COLOR;
        vec3 specular = SPECULAR_INTENSITY * (spec + fresnel * 0.5) * LIGHT_COLOR;
        vec3 rimLight = rim * LIGHT_COLOR;
        
        vec3 color = ambient + diffuse + specular + rimLight;
        
        // Improved tone mapping (ACES-like)
        color = color * (2.51 * color + 0.03) / (color * (2.43 * color + 0.59) + 0.14);
        
        fragColor = vec4(color, 1.0);
    } else {
        // Background color from skybox
        fragColor = texture(iChannel3, uv);
    }
}
